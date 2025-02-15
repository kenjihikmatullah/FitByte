package activityService

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	authJwt "github.com/rafitanujaya/go-fiber-template/src/auth/jwt"
	functionCallerInfo "github.com/rafitanujaya/go-fiber-template/src/logger/helper"
	loggerZap "github.com/rafitanujaya/go-fiber-template/src/logger/zap"
	"github.com/rafitanujaya/go-fiber-template/src/model/dtos/request"
	"github.com/rafitanujaya/go-fiber-template/src/model/dtos/response"
	Entity "github.com/rafitanujaya/go-fiber-template/src/model/entities/activity"
	activityRepository "github.com/rafitanujaya/go-fiber-template/src/repositories/activity"
	"github.com/samber/do/v2"
)

type activityService struct {
	ActivityRepository activityRepository.ActivityRepositoryInterface
	Db                 *pgxpool.Pool
	jwtService         authJwt.JwtServiceInterface
	Logger             loggerZap.LoggerInterface
}

func NewActivityService(activityRepo activityRepository.ActivityRepositoryInterface, db *pgxpool.Pool, jwtService authJwt.JwtServiceInterface, logger loggerZap.LoggerInterface) ActivityServiceInterface {
	return &activityService{ActivityRepository: activityRepo, Db: db, jwtService: jwtService, Logger: logger}
}

func NewActivityServiceInject(i do.Injector) (ActivityServiceInterface, error) {
	_db := do.MustInvoke[*pgxpool.Pool](i)
	_activityRepo := do.MustInvoke[activityRepository.ActivityRepositoryInterface](i)
	_jwtService := do.MustInvoke[authJwt.JwtServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	return NewActivityService(_activityRepo, _db, _jwtService, _logger), nil
}

func (as *activityService) Create(ctx context.Context, input request.RequestActivity) (response.ResponseActivity, error) {
	activity := Entity.Activity{}

	timeNow := time.Now()
	activity.CreatedAt = timeNow
	activity.UpdatedAt = timeNow
	activity.ActivityType = (Entity.ActivityType)(*input.ActivityType)
	activity.UserId = *input.UserId
	activity.DurationInMinutes = int64(*input.DurationInMinutes)
	activity.CaloriesBurned = Entity.CountCalories(activity.DurationInMinutes, activity.ActivityType)

	doneAt, err := time.Parse(time.RFC3339, *input.DoneAt)
	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceCreate)
		return response.ResponseActivity{}, err
	}
	activity.DoneAt = doneAt

	activity.ActivityId, err = as.ActivityRepository.Create(ctx, as.Db, activity)

	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceCreate)
		return response.ResponseActivity{}, err
	}

	return response.ResponseActivity{
		ActivityId:        activity.ActivityId,
		ActivityType:      *input.ActivityType,
		DoneAt:            activity.DoneAt.Format(time.RFC3339),
		DurationInMinutes: int(activity.DurationInMinutes),
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         activity.UpdatedAt.Format(time.RFC3339),
	}, nil
}
