package clean

import (
	"context"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/util/runtime"
	"time"

	"github.com/horizoncd/horizon/core/common"
	"github.com/horizoncd/horizon/lib/q"
	"github.com/horizoncd/horizon/pkg/config/clean"
	"github.com/horizoncd/horizon/pkg/event/models"
	"github.com/horizoncd/horizon/pkg/param/managerparam"
	"github.com/horizoncd/horizon/pkg/util/log"
)

type Cleaner struct {
	clean.Config
	eventRules       map[string][]clean.EventCleanRule
	mgr              *managerparam.Manager
	eventCursor      uint
	webhookLogCursor uint
}

func New(config clean.Config, mgr *managerparam.Manager) *Cleaner {
	if config.Batch == 0 {
		config.Batch = 160
	}
	eventRules := make(map[string][]clean.EventCleanRule, len(config.EventCleanRules))
	for _, rule := range config.EventCleanRules {
		eventRules[rule.EventType] = append(eventRules[rule.EventType], rule)
	}
	return &Cleaner{
		Config:     config,
		eventRules: eventRules,
		mgr:        mgr,
	}
}

func (c *Cleaner) Run(ctx context.Context) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	cron := cron.New(cron.WithSeconds(), cron.WithLocation(loc))
	_, err = cron.AddFunc(c.TimeToRun, func() {
		defer runtime.HandleCrash()
		log.Debugf(ctx, "start to clean")
		c.webhookLogClean(ctx)
		c.eventClean(ctx)
	})
	if err != nil {
		panic(err)
	}
	cron.Run()
}

func (c *Cleaner) webhookLogClean(ctx context.Context) {
	log.Debugf(ctx, "start to clean webhooklogs")
	defer func() {
		c.webhookLogCursor = 0
		log.Debugf(ctx, "finish to clean webhooklogs")
	}()
	if len(c.WebhookLogCleanRules) == 0 {
		return
	}
	needDeleted := make([]uint, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		query := &q.Query{
			Keywords: q.KeyWords{
				common.StartID: c.webhookLogCursor,
				common.Limit:   c.Batch,
				common.OrderBy: "l.id",
			},
		}
		webhooklogs, _, err := c.mgr.WebhookManager.ListWebhookLogs(ctx, query, nil)
		if err != nil {
			log.Errorf(ctx, "failed to list webhooklogs: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(webhooklogs) == 0 {
			return
		}

		maxID := uint(0)
		needDeleted = needDeleted[:0]
		for _, webhooklog := range webhooklogs {
			if webhooklog.ID > maxID {
				maxID = webhooklog.ID
			}
			event, err := c.mgr.EventManager.GetEvent(ctx, webhooklog.EventID)
			if err != nil {
				log.Errorf(ctx, "failed to get event: %v", err)
				continue
			}
			for _, rule := range c.WebhookLogCleanRules {
				if webhooklog.UpdatedAt.Add(rule.TTL).After(time.Now()) {
					continue
				}
				if rule.RelatedEventType != event.EventType {
					continue
				}
				needDeleted = append(needDeleted, webhooklog.ID)
			}
		}

		c.webhookLogCursor = maxID
		if len(needDeleted) == 0 {
			continue
		}

		log.Debugf(ctx, "need to delete webhooklogs: %v", needDeleted)
		_, _ = c.mgr.WebhookManager.DeleteWebhookLogs(ctx, needDeleted...)
	}
}

func (c *Cleaner) needClean(ctx context.Context, event *models.Event) bool {
	rules := c.eventRules[event.EventType]
	for _, rule := range rules {
		if time.Now().Before(event.CreatedAt.Add(rule.TTL)) {
			continue
		}
		m := make(map[string]interface{})
		if event.Extra != nil && *event.Extra != "" {
			err := json.Unmarshal([]byte(*event.Extra), &m)
			if err != nil {
				log.Errorf(ctx, "failed to unmarshal event extra: %v", err)
				return true
			}
			if rule.Reason != "" && rule.Reason != m["reason"] {
				continue
			}
			involvedObject := m["involvedObject"].(map[string]interface{})
			if involvedObject != nil {
				if rule.APIVersion != "" && rule.APIVersion != involvedObject["apiVersion"] {
					continue
				}
				if rule.Kind != "" && rule.Kind != involvedObject["kind"] {
					continue
				}
				if rule.Name != "" && rule.Name != involvedObject["name"] {
					continue
				}
				if rule.Namespace != "" && rule.Namespace != involvedObject["namespace"] {
					continue
				}
			}
		}
		return true
	}
	return false
}

func (c *Cleaner) eventClean(ctx context.Context) {
	log.Debugf(ctx, "start to clean events")
	defer func() {
		c.eventCursor = 0
		log.Debugf(ctx, "finish to clean events")
	}()
	if len(c.EventCleanRules) == 0 {
		return
	}
	needDeleted := make([]uint, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		events, err := c.mgr.EventManager.ListEvents(ctx, &q.Query{Keywords: q.KeyWords{
			common.Limit:   c.Batch,
			common.StartID: int(c.eventCursor),
		}})
		if err != nil {
			log.Errorf(ctx, "failed to list events: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if len(events) == 0 {
			return
		}

		maxID := uint(0)
		needDeleted = needDeleted[:0]
		for _, event := range events {
			if event.ID > maxID {
				maxID = event.ID
			}
			if c.needClean(ctx, event) {
				needDeleted = append(needDeleted, event.ID)
			}
		}
		c.eventCursor = maxID
		if len(needDeleted) == 0 {
			continue
		}
		log.Infof(ctx, "need to delete event: %v", needDeleted)
		_, _ = c.mgr.EventManager.DeleteEvents(ctx, needDeleted...)
	}
}
