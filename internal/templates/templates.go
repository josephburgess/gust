package templates

import (
	"html/template"
	"io"
	"time"
)

const authSuccessTemplateContent = `<!DOCTYPE html>
<html>
<head>
    <title>Gust Authentication Success</title>
    <style>
        :root {
            --base: #191724;
            --surface: #1f1d2e;
            --overlay: #26233a;
            --subtle: #908caa;
            --text: #e0def4;
            --love: #eb6f92;
            --gold: #f6c177;
            --rose: #ebbcba;
            --pine: #31748f;
            --foam: #9ccfd8;
            --iris: #c4a7e7;
            --highlight-med: #403d52;
        }
        body {
            font-family: system-ui, sans-serif;
            max-width: 600px;
            margin: 0 auto;
            padding: 2rem;
            text-align: center;
            line-height: 1.6;
            background-color: var(--base);
            color: var(--text);
        }
        h1 {
            color: var(--rose);
            margin-bottom: 1.5rem;
        }
        .success {
            color: var(--pine);
            font-weight: bold;
            font-size: 1.2rem;
        }
        .info {
            margin: 2rem 0;
            line-height: 1.5;
            color: var(--subtle);
        }
        .api-key {
            background: var(--surface);
            padding: 1rem;
            border-radius: 4px;
            font-family: monospace;
            overflow-wrap: break-word;
            margin: 1.5rem 0;
            text-align: center;
            border: 1px solid var(--highlight-med);
            color: var(--foam);
        }
        .code-block {
            background: var(--overlay);
            color: var(--text);
            padding: 1rem;
            border-radius: 4px;
            font-family: monospace;
            overflow-wrap: break-word;
            margin: 1.5rem 0;
            text-align: left;
            white-space: pre;
            border: 1px solid var(--highlight-med);
        }
        .next-steps {
            background: var(--surface);
            padding: 1.5rem;
            border-radius: 6px;
            margin-top: 2rem;
            border: 1px solid var(--highlight-med);
        }
        .command {
            background: var(--overlay);
            color: var(--foam);
            padding: 0.5rem 1rem;
            border-radius: 4px;
            font-family: monospace;
            display: inline-block;
            margin: 0.5rem 0;
        }
        .key-highlight {
            color: var(--love);
        }
        .value-highlight {
            color: var(--gold);
        }
        .string-highlight {
            color: var(--iris);
        }
    </style>
</head>
<body>
    <h1>Authentication Successful!</h1>
    <p class="success">Welcome, {{.Login}}!</p>
    <p class="info">Your Gust API key has been generated:</p>
    <div class="api-key">{{.ApiKey}}</div>
    <p>This key will allow you to access weather data through the breeze API.</p>
    <div class="next-steps">
        <p>You can now return to your terminal. The CLI application should automatically continue.</p>
        <p>If it doesn't, you can close this window and add the below to your ~/.config/gust/auth.json</p>
        <div class="code-block">{
  <span class="key-highlight">"api_key"</span>: <span class="string-highlight">"{{.ApiKey}}"</span>,
  <span class="key-highlight">"server_url"</span>: <span class="string-highlight">"{{.ServerURL}}"</span>,
  <span class="key-highlight">"github_user"</span>: <span class="string-highlight">"{{.Login}}"</span>,
  <span class="key-highlight">"last_auth"</span>: <span class="string-highlight">"{{.LastAuth}}"</span>
}</div>
    </div>
</body>
</html>`

var templates *template.Template

func init() {
	templates = template.Must(template.New("auth_success").Parse(authSuccessTemplateContent))
}

func RenderSuccessTemplate(w io.Writer, login, apiKey, serverURL string) error {
	lastAuth := time.Now().Format(time.RFC3339)

	data := struct {
		Login     string
		ApiKey    string
		ServerURL string
		LastAuth  string
	}{
		Login:     login,
		ApiKey:    apiKey,
		ServerURL: serverURL,
		LastAuth:  lastAuth,
	}

	return templates.ExecuteTemplate(w, "auth_success", data)
}
