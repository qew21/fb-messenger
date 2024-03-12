package assistant

import (
	"fmt"
	"testing"

	"github.com/qew21/fb-messenger/config"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestQianWen(t *testing.T) {
	userID := "test_user"
	message := "Where can I find this film?"
	appConfig, _ := config.LoadConfig("../config.yaml")
	_, err := QianWen(userID, message, appConfig.QianwenKey)
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf("Error reply for message '%s'", message))
		return
	}
	assert.Equal(t, nil, err, "Error reply for message '%s'", message)
}
