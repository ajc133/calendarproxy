package handlers

import (
	"log"
	"net/http"

	"github.com/ajc133/calendarproxy/pkg/calendar"
	"github.com/ajc133/calendarproxy/pkg/db"
	"github.com/gin-gonic/gin"
)

const ContentType string = "Content-Type"
const CalendarContent string = "text/calendar; charset=utf-8"

type CreateParams struct {
	Url                string `form:"url" json:"url" binding:"required"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary" binding:"required"`
}

type UpdateParams struct {
	Id                 string `form:"id" json:"id" binding:"required"`
	Url                string `form:"url" json:"url"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary"`
}

func CreateCalendar(c *gin.Context) {
	json := CreateParams{}
	err := c.Bind(&json)
	if err != nil {
		log.Printf("Error: %s", err)

		// TODO: User-friendly error for form submission
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Parameters submitted were malformed"})
		return
	}

	record := db.Record{
		Url:                json.Url,
		ReplacementSummary: json.ReplacementSummary,
	}

	id, err := db.WriteRecord(record)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error writing to db"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

func UpdateCalendar(c *gin.Context) {
	json := UpdateParams{}
	err := c.Bind(&json)
	if err != nil {
		log.Printf("Error: %s", err)

		// TODO: User-friendly error for form submission
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Parameters submitted were malformed"})
		return
	}

	record := db.ChangeRecord{
		Id:                 json.Id,
		Url:                json.Url,
		ReplacementSummary: json.ReplacementSummary,
	}

	_, err = db.UpdateRecord(record)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error writing to db"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})

}

func GetCalendarByID(c *gin.Context) {
	id := c.Param("id")
	record, err := db.ReadRecord(id)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error in database lookup"})
		return
	}

	newCal, err := calendar.FetchAndTransformCalendar(record.Url, record.ReplacementSummary)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when fetching and parsing the given URL"})
	}
	c.Header(ContentType, CalendarContent)
	c.String(http.StatusOK, newCal)
}
