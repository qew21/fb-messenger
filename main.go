/**
 * This source code is licensed under the license found in the
 * LICENSE file in the root directory of this source tree.
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qew21/fb-messenger/config"
	"github.com/qew21/fb-messenger/messenger"

	"github.com/julienschmidt/httprouter"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	receivedUpdatesMutex sync.RWMutex
	receivedUpdates      = make([]map[string]interface{}, 0)
	appConfig            *config.AppConfig
)

func main() {
	environment := strings.ToLower(os.Getenv("ENVIRONMENT"))
	var logFilePath = "log/app.log"

	if strings.EqualFold(environment, "development") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		err := os.MkdirAll(filepath.Dir(logFilePath), 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create log directory")
		}

		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open or create log file")
		}
		defer file.Close()
		log.Logger = log.Output(zerolog.New(file).With().Timestamp().Logger())
	}
	var err error
	appConfig, err = config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading configuration")
	}

	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", handleGetIndex)
	router.HandlerFunc(http.MethodGet, "/facebook", handleGetWebHook)
	router.HandlerFunc(http.MethodPost, "/facebook", handlePostWebHook)
	router.HandlerFunc(http.MethodGet, "/privacy", privacyHandler)
	router.HandlerFunc(http.MethodGet, "/terms", termsHandler)

	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	log.Info().Str("address", addr).Msg("Starting HTTP server")
	if appConfig.Port == 443 {
		err = http.ListenAndServeTLS(addr, appConfig.CertFile, appConfig.KeyFile, router)
	} else {
		err = http.ListenAndServe(addr, router)
	}
	if err != nil {
		log.Fatal().Str("address", addr).Err(err).Msg("Error starting HTTP server")
	}

}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	receivedUpdatesMutex.RLock()
	defer receivedUpdatesMutex.RUnlock()
	json.NewEncoder(w).Encode(receivedUpdates)
}

func handleGetWebHook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("hub.mode") != "subscribe" || query.Get("hub.verify_token") != appConfig.Token {
		log.Warn().Msg("Invalid subscribe token")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	challengeNum, err := strconv.Atoi(query.Get("hub.challenge"))
	if err != nil {
		log.Warn().Err(err).Msg("Error converting hub.challenge to integer")
		http.Error(w, "Invalid challenge value", http.StatusBadRequest)
		return
	}

	w.Write([]byte(strconv.Itoa(challengeNum)))
}

func handlePostWebHook(w http.ResponseWriter, r *http.Request) {
	payloadBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Error reading request body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding payload")
		handleJSONUnmarshalError(w, err)
		return
	}

	err = messenger.ProcessMessage(payload, appConfig, r.URL.Path == "/test")
	if err != nil {
		log.Warn().Err(err).Str("object", payload["object"].(string)).Msg("ProcessMessage Failed")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook processed successfully.")
}

func handleJSONUnmarshalError(w http.ResponseWriter, err error) {
	if serr, ok := err.(*json.SyntaxError); ok {
		http.Error(w, serr.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func termsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	termsHTML := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>User Terms of Service</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 40px; }
			h1 { color: #333366; }
			p { margin: 20px 0; }
			ul { margin: 20px 0; }
		</style>
	</head>
	<body>
		<h1>User Terms of Service</h1>
		<p>Welcome to his application! By using our application, you agree to the following terms and conditions:</p>
		<ul>
			<li><strong>Acceptance of Terms</strong>: When you access our application, you agree to be bound by these Terms of Service.</li>
			<li><strong>Modification of Terms</strong>: We reserve the right to modify these terms at any time. Your continued use of the application signifies your acceptance of any adjustments.</li>
			<li><strong>User Conduct</strong>: You are responsible for all your activity in connection with the service and ensuring that all content uploaded complies with applicable laws and regulations.</li>
			<li><strong>Intellectual Property</strong>: All content included on the application, such as text, graphics, logos, and software, is the property of his application or its content suppliers.</li>
		</ul>
		<p>For more information or if you have any questions, please contact us at this site.</p>
	</body>
	</html>
	`
	fmt.Fprint(w, termsHTML)
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	privacyHTML := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Privacy Policy</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 40px; }
			h1 { color: #333366; }
			p { margin: 20px 0; }
			ul { margin: 20px 0; }
		</style>
	</head>
	<body>
		<h1>Privacy Policy</h1>
		<p>At this application, we respect the privacy of our users. This Privacy Policy outlines how we collect, use, and protect your personal information:</p>
		<ul>
			<li><strong>Data Collection</strong>: We collect information such as your name, email address, and usage data to provide and improve our service.</li>
			<li><strong>Use of Information</strong>: Your information helps us to personalize the service and improve your experience.</li>
			<li><strong>Sharing of Information</strong>: We do not sell, trade, or otherwise transfer to outside parties your personally identifiable information without your consent.</li>
			<li><strong>Data Security</strong>: We implement a variety of security measures to maintain the safety of your personal information.</li>
		</ul>
		<p>We may update our Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page.</p>
		<p>If you have any questions about this Privacy Policy, please contact us at this site.</p>
	</body>
	</html>
	`
	fmt.Fprint(w, privacyHTML)
}
