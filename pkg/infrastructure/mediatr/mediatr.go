package mediatr

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/goccy/go-reflect"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/backend-api/pkg/constant"
	res "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/logger"
	"github.com/ndodanli/backend-api/pkg/servers/lifetime"
	"github.com/pkg/errors"
)

// RequestHandlerFunc is a continuation for the next task to execute in the pipeline
type RequestHandlerFunc func(echoCtx echo.Context) interface{}

// PipelineBehavior is a Pipeline behavior for wrapping the inner handler.
type PipelineBehavior interface {
	Handle(echoCtx echo.Context, request interface{}, next RequestHandlerFunc) interface{}
}

type RequestHandler[TRequest any, TResponse any] interface {
	Handle(echoCtx echo.Context, request TRequest) TResponse
}

type RequestHandlerFactory[TRequest any, TResponse any] func() RequestHandler[TRequest, TResponse]

type NotificationHandler[TNotification any] interface {
	Handle(echoCtx echo.Context, notification TNotification) error
}

type NotificationHandlerFactory[TNotification any] func() NotificationHandler[TNotification]

var requestHandlersRegistrations = map[reflect.Type]interface{}{}
var notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
var pipelineBehaviours []interface{} = []interface{}{}

type Unit struct{}

// RegisterRequestHandler register the request handler to mediatr registry.
func RegisterRequestHandler[TRequest any, TResponse any](handler RequestHandler[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](handler)
}

// RegisterRequestHandlerFactory register the request handler factory to mediatr registry.
func RegisterRequestHandlerFactory[TRequest any, TResponse any](factory RequestHandlerFactory[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](factory)
}

// RegisterRequestPipelineBehaviors register the request behaviors to mediatr registry.
func RegisterRequestPipelineBehaviors(behaviours ...PipelineBehavior) error {
	for _, behavior := range behaviours {
		behaviorType := reflect.TypeOf(behavior)

		existsPipe := existsPipeType(behaviorType)
		if existsPipe {
			return errors.Errorf("registered behavior already exists in the registry.")
		}

		pipelineBehaviours = append(pipelineBehaviours, behavior)
	}

	return nil
}

func registerNotificationHandler[TEvent any](handler any) error {
	var event TEvent
	eventType := reflect.TypeOf(event)

	handlers, exist := notificationHandlersRegistrations[eventType]
	if !exist {
		notificationHandlersRegistrations[eventType] = []interface{}{handler}
		return nil
	}

	notificationHandlersRegistrations[eventType] = append(handlers, handler)

	return nil
}

// RegisterNotificationHandler register the notification handler to mediatr registry.
func RegisterNotificationHandler[TEvent any](handler NotificationHandler[TEvent]) error {
	return registerNotificationHandler[TEvent](handler)
}

// RegisterNotificationHandlerFactory register the notification handler factory to mediatr registry.
func RegisterNotificationHandlerFactory[TEvent any](factory NotificationHandlerFactory[TEvent]) error {
	return registerNotificationHandler[TEvent](factory)
}

// RegisterNotificationHandlers register the notification handlers to mediatr registry.
func RegisterNotificationHandlers[TEvent any](handlers ...NotificationHandler[TEvent]) error {
	if len(handlers) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range handlers {
		err := RegisterNotificationHandler(handler)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterNotificationHandlersFactories register the notification handlers factories to mediatr registry.
func RegisterNotificationHandlersFactories[TEvent any](factories ...NotificationHandlerFactory[TEvent]) error {
	if len(factories) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range factories {
		err := RegisterNotificationHandlerFactory[TEvent](handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearRequestRegistrations() {
	requestHandlersRegistrations = map[reflect.Type]interface{}{}
}

func ClearNotificationRegistrations() {
	notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
}

func buildRequestHandler[TRequest any, TResponse any](handler any) (RequestHandler[TRequest, TResponse], bool) {
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		factory, ok := handler.(RequestHandlerFactory[TRequest, TResponse])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

func registerRequestHandler[TRequest any, TResponse any](handler any) error {
	var request TRequest
	requestType := reflect.TypeOf(request)

	_, exist := requestHandlersRegistrations[requestType]
	if exist {
		// each request in request/response strategy should have just one handler
		//log and return nil
		fmt.Printf("registered handler already exists in the registry for message %s\n", requestType.String())
		return nil
	}

	handleValElem := reflect.ValueOf(handler).Elem()
	handlerType := handleValElem.Type()

	for i := 0; i < handleValElem.NumField(); i++ {
		field := handleValElem.Field(i)
		fieldType := handlerType.Field(i).Type

		if field.IsValid() && field.IsNil() {
			if fieldType.Kind() == reflect.Interface {
				if fieldType.Implements(reflect.TypeOf((*logger.ILogger)(nil)).Elem()) {
					field.Set(reflect.ValueOf(lifetime.LoggerSingleton))
				} else if fieldType.Implements(reflect.TypeOf((*uow.IUnitOfWork)(nil)).Elem()) {
					field.Set(reflect.ValueOf(lifetime.UOWSingleton))
				}
			}
			switch fieldType {
			case reflect.TypeOf(lifetime.AppServicesType):
				field.Set(reflect.ValueOf(lifetime.AppServicesSingleton))
			}
		}
	}

	requestHandlersRegistrations[requestType] = handler

	return nil
}

// Send the request to its corresponding request handler.
func Send[TRequest any, TResponse any](echoCtx echo.Context, request TRequest) TResponse {
	requestType := reflect.TypeOf(request)
	var response TResponse
	handler, ok := requestHandlersRegistrations[requestType]
	if !ok {
		// request-response strategy should have exactly one handler and if we can't find a corresponding handler, we should return an error
		//return *new(TResponse), errors.Errorf("no handler for request %T", request)
		echoCtx.Logger().Error("no handler for request %T", request)
		baseHttpApiResult := res.NewResult[any, *echo.HTTPError, any](nil)
		baseHttpApiResult.SetErrorMessage("Internal Server Error")
		echoCtx.JSON(500, baseHttpApiResult)
		return response
	}

	handleValElem := reflect.ValueOf(handler).Elem()
	handlerType := handleValElem.Type()

	// Handle scoped instances
	for i := 0; i < handleValElem.NumField(); i++ {
		field := handleValElem.Field(i)
		fieldType := handlerType.Field(i).Type

		switch fieldType {
		case reflect.TypeOf(lifetime.TxSessionManagerType):
			field.Set(reflect.ValueOf(echoCtx.Get(constant.General.TxSessionManagerKey).(*pg.TxSessionManager)))
		}
	}

	handlerValue, ok := buildRequestHandler[TRequest, TResponse](handler)
	if !ok {
		//return *new(TResponse), errors.Errorf("handler for request %T is not a Handler", request)
		echoCtx.Logger().Error("handler for request %T is not a Handler", request)
		baseHttpApiResult := res.NewResult[any, *echo.HTTPError, any](nil)
		baseHttpApiResult.SetErrorMessage("Internal Server Error")
		echoCtx.JSON(500, baseHttpApiResult)
		return response
	}

	if len(pipelineBehaviours) > 0 {
		var reversPipes = reversOrder(pipelineBehaviours)

		var lastHandler RequestHandlerFunc = func(echoCtx echo.Context) interface{} {
			return handlerValue.Handle(echoCtx, request)
		}

		aggregateResult := linq.From(reversPipes).AggregateWithSeedT(lastHandler, func(next RequestHandlerFunc, pipe PipelineBehavior) RequestHandlerFunc {
			pipeValue := pipe
			nexValue := next

			var handlerFunc RequestHandlerFunc = func(echoCtx echo.Context) interface{} {
				return pipeValue.Handle(echoCtx, request, nexValue)
			}

			return handlerFunc
		})

		v := aggregateResult.(RequestHandlerFunc)
		response := v(echoCtx)

		//if err != nil {
		//	return *new(TResponse), errors.Wrap(err, "error handling request")
		//}

		return response.(TResponse)
	} else {
		res := handlerValue.Handle(echoCtx, request)
		//if err != nil {
		//	return *new(TResponse), errors.Wrap(err, "error handling request")
		//}

		response = res
	}

	return response
}

func buildNotificationHandler[TNotification any](handler any) (NotificationHandler[TNotification], bool) {
	handlerValue, ok := handler.(NotificationHandler[TNotification])
	if !ok {
		factory, ok := handler.(NotificationHandlerFactory[TNotification])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

// Publish the notification event to its corresponding notification handler.
func Publish[TNotification any](echoCtx echo.Context, notification TNotification) error {
	eventType := reflect.TypeOf(notification)

	handlers, ok := notificationHandlersRegistrations[eventType]
	if !ok {
		// notification strategy should have zero or more handlers, so it should run without any error if we can't find a corresponding handler
		return nil
	}

	for _, handler := range handlers {
		handlerValue, ok := buildNotificationHandler[TNotification](handler)

		if !ok {
			return errors.Errorf("handler for notification %T is not a Handler", notification)
		}

		err := handlerValue.Handle(echoCtx, notification)
		if err != nil {
			return errors.Wrap(err, "error handling notification")
		}
	}

	return nil
}

func reversOrder(values []interface{}) []interface{} {
	var reverseValues []interface{}

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}

func existsPipeType(p reflect.Type) bool {
	for _, pipe := range pipelineBehaviours {
		if reflect.TypeOf(pipe) == p {
			return true
		}
	}

	return false
}
