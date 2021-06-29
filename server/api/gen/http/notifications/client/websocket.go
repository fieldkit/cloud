// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notifications WebSocket client streaming
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"io"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"
	notificationsviews "github.com/fieldkit/cloud/server/api/gen/notifications/views"
	"github.com/gorilla/websocket"
	goahttp "goa.design/goa/v3/http"
)

// ConnConfigurer holds the websocket connection configurer functions for the
// streaming endpoints in "notifications" service.
type ConnConfigurer struct {
	ListenFn goahttp.ConnConfigureFunc
}

// ListenClientStream implements the notifications.ListenClientStream interface.
type ListenClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NewConnConfigurer initializes the websocket connection configurer function
// with fn for all the streaming endpoints in "notifications" service.
func NewConnConfigurer(fn goahttp.ConnConfigureFunc) *ConnConfigurer {
	return &ConnConfigurer{
		ListenFn: fn,
	}
}

// Recv reads instances of "notifications.Notification" from the "listen"
// endpoint websocket connection.
func (s *ListenClientStream) Recv() (*notifications.Notification, error) {
	var (
		rv   *notifications.Notification
		body ListenResponseBody
		err  error
	)
	err = s.conn.ReadJSON(&body)
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		s.conn.Close()
		return rv, io.EOF
	}
	if err != nil {
		return rv, err
	}
	res := NewListenNotificationOK(&body)
	vres := &notificationsviews.Notification{res, "default"}
	if err := notificationsviews.ValidateNotification(vres); err != nil {
		return rv, goahttp.ErrValidationError("notifications", "listen", err)
	}
	return notifications.NewNotification(vres), nil
}
