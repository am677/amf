/*
 * Namf_EventExposure
 *
 * AMF Event Exposure Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package eventexposure

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/context"
	"free5gc/src/amf/producer"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// CreateSubscription - Namf_EventExposure Subscribe service Operation
func CreateSubscription(c *gin.Context) {

	var createEventSubscription models.AmfCreateEventSubscription

	if err := c.ShouldBindJSON(&createEventSubscription); err != nil {
		log.Panic(err.Error())
	}
	self := context.AMF_Self()
	res, problem := producer.CreateAMFEventSubscription(self, createEventSubscription, time.Now().UTC())
	if problem.Cause != "" {
		c.JSON(int(problem.Status), problem)
	} else {
		c.JSON(http.StatusCreated, res)
	}
}
