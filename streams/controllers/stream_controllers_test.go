package controllers_test

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ello/ello-go/streams/controllers"
	"github.com/m4rw3r/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamController", func() {
	var id uuid.UUID

	BeforeEach(func() {
		id, _ = uuid.V4()
	})

	Context("when adding content via PUT /streams", func() {

		It("should return a status 201 when passed a correct body", func() {
			item1ID, _ := uuid.V4()
			item2ID, _ := uuid.V4()
			items := []controllers.StreamItem{{
				StreamID:  id,
				Timestamp: time.Now(),
				Type:      0,
				ID:        item1ID,
			}, {
				StreamID:  id,
				Timestamp: time.Now(),
				Type:      1,
				ID:        item2ID,
			}}
			itemsJSON, _ := json.Marshal(items)
			Request("PUT", "/streams", string(itemsJSON))
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))
			//TODO Verify it tries to add to the StreamService
		})

		It("should return a status 201 when passed a correct body string", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				},
				{
					"id":"c8f17401-62d0-444c-a5d6-639b01f6070f",
					"ts":"2015-11-16T11:59:29.313068877-07:00",
					"type":1,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		It("should return a status 422 when passed an invalid uuid", func() {
			jsonStr := `[
				{
					"id":"ABC",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid date (non ISO8601)", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
		It("should return a status 422 when passed an invalid type", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":a,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when validation error is in later element", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				},
				{
					"id":"c8f17401-62d0-444c-a5d6-639b01f6070f",
					"ts":"2015-11-16T11:59:29.313068877-07:00",
					"type":a,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid body/query", func() {
			Request("PUT", "/streams", "hi")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
	})
	Context("when retrieving a stream via /stream/:id", func() {

		It("should return a status 201 when accessed with a valid ID", func() {
			Request("GET", "/stream/"+id.String(), "")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 422 when passed an invalid id", func() {
			Request("GET", "/stream/"+"abc123", "")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
	})
	Context("when retrieving streams via /streams/coalesce", func() {

		It("should return a status 201 when accessed with a valid ID", func() {
			Request("POST", "/streams/coalesce", "")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 422 when passed an invalid query", func() {
			Request("POST", "/streams/coalesce", "")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		// 	It("should return a status 200 with no args", func() {
		// 		Request("GET", "/users")
		// 		var data []service.User
		// 		_ = json.Unmarshal(response.Body.Bytes(), &data)
		//
		// 		Expect(response.Code).To(Equal(http.StatusOK))
		// 		Expect(data[0].Username).To(Equal("rtyer"))
		// 	})
		//
		// 	It("should use the passed limit/offset", func() {
		// 		Request("GET", "/users?limit=5&offset=13")
		//
		// 		Expect(response.Code).To(Equal(http.StatusOK))
		// 		Expect(userService.lastLimit).To(Equal(5))
		// 		Expect(userService.lastOffset).To(Equal(13))
		// 	})
		//
		// 	It("should correctly validate the limit", func() {
		// 		Request("GET", "/users?limit=a")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotAcceptable))
		// 	})
		//
		// 	It("should correctly validate the offset", func() {
		// 		Request("GET", "/users?offset=a")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotAcceptable))
		// 	})
		// })
		// Context("when calling /user/<username>", func() {
		//
		// 	It("should return a status 200 with a user that is present", func() {
		// 		Request("GET", "/users/rtyer")
		// 		var user service.User
		// 		_ = json.Unmarshal(response.Body.Bytes(), &user)
		//
		// 		Expect(response.Code).To(Equal(http.StatusOK))
		// 		Expect(user.Username).To(Equal("rtyer"))
		// 	})
		//
		// 	It("should return a status 404 with a non existent user", func() {
		// 		Request("GET", "/users/asdf")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotFound))
		// 	})
		//
		// 	It("should return a status 406 if the username is invalid", func() {
		// 		Request("GET", "/users/^&*$")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotAcceptable))
		// 	})
		//
		// 	It("should accept UTF-8 characters for the username", func() {
		// 		Request("GET", "/users/ßåœ")
		// 		Expect(response.Code).NotTo(Equal(http.StatusNotAcceptable))
		// 	})
	})
})