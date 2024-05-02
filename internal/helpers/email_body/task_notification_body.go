package emailbody

import (
	"bytes"
	"html/template"
	"log"
	"todoapp/internal/database/model/todo"
)

func Generate_task_notification_email_body(username string, taskDetails todo.Task) []byte {
	data := struct {
		TaskName    string
		DueDate     string
		UserName    string
		CompanyName string
	}{
		TaskName:    taskDetails.Title,
		DueDate:     taskDetails.Duedate.Format("2006-01-02"),
		UserName:    username,
		CompanyName: "TodoApp"}

	const emailTemplate = `
				
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>Task Due is: {{.TaskName}}</title>
				</head>
				<body>
					<p>Dear {{.UserName}},</p>
					<p>This is a reminder that the following task is due soon:</p>
					<ul>
						<li><strong>Task Name:</strong> {{.TaskName}}</li>
						<li><strong>Task Due Date:</strong> {{.DueDate}}</li>
					</ul>
					<p>Please ensure that you complete this task before the due date to avoid any delays.</p>
					<p>Thank you,<br>{{.CompanyName}}</p>
				</body>
				</html>
			`

	// Parse the HTML template
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the template with the data
	var emailBody bytes.Buffer
	err = tmpl.Execute(&emailBody, data)
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte(
		"Subject: Task Due: " + data.TaskName + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" + emailBody.String())

	return msg
}
