package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
)

var (
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
)

func InitEmail() {
	smtpHost = os.Getenv("SMTPHOST")
	smtpPort = os.Getenv("SMTPPORT")
	smtpUser = os.Getenv("SMTPUSER")
	smtpPass = os.Getenv("SMTPPASS")
}

func SendEmailToVerify(userEmail string, code string) {
	println(smtpUser)
	println("STARTING TO SEND")
	e := email.NewEmail()
	e.From = fmt.Sprintf("Formular - школа по подготовке к математике! <%s>", smtpUser)
	e.To = []string{userEmail}
	e.Subject = "Подтверждение вашей электронной почты"

	data := struct {
		Name string
		User string
		Date string
		Code string
	}{
		Name: "Ученик",
		User: userEmail,
		Date: time.Now().Format("02 January 2006"),
		Code: code,
	}

	tmpl := template.Must(template.New("verify").Parse(VERIFYHTMLTEMPLATE))

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("Ошибка рендеринга HTML: %v", err)
		return
	}

	e.HTML = htmlBody.Bytes()
	e.Text = []byte(fmt.Sprintf(
		"Подтверждение email для Formular\n\n"+
			"Здравствуйте!\n\n"+
			"Ваш код подтверждения: %s\n\n"+
			"Введите этот код в приложении для завершения регистрации.\n\n"+
			"Если вы не запрашивали это подтверждение, проигнорируйте данное письмо.\n\n"+
			"С уважением,\nКоманда Formular",
		code))

	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	tlsconfig := &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: true,
	}
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	if err := e.SendWithTLS(
		smtpAddr,
		auth,
		tlsconfig,
	); err != nil {
		log.Printf("Ошибка отправки письма: %v", err)
	} else {
		log.Printf("Письмо с подтверждением отправлено на %s", userEmail)
	}
}

const VERIFYHTMLTEMPLATE = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Подтверждение электронной почты</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f5f7fa; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; color: #333333; line-height: 1.6;">
    <table width="100%" cellpadding="0" cellspacing="0" style="max-width: 650px; margin: 30px auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 5px 25px rgba(0, 0, 0, 0.08);">
        <!-- Header with brand identity -->
        <tr>
            <td style="background: linear-gradient(135deg, #4361ee 0%, #3a0ca3 100%); padding: 35px 30px; text-align: center; color: white;">
                <div style="font-size: 28px; font-weight: 700; letter-spacing: -0.5px; margin-bottom: 8px;">Formular</div>
                <div style="font-size: 18px; opacity: 0.9; font-weight: 300;">Школа подготовки к математике</div>
                <div style="width: 60px; height: 4px; background: rgba(255, 255, 255, 0.3); margin: 15px auto; border-radius: 2px;"></div>
            </td>
        </tr>
        
        <!-- Main content -->
        <tr>
            <td style="padding: 40px 30px;">
                <h1 style="color: #2b2d42; font-size: 24px; margin-top: 0; margin-bottom: 25px;">Здравствуйте, {{.Name}}!</h1>
                
                <p style="font-size: 16px; margin-bottom: 30px; color: #495057;">
                    Спасибо, что присоединились к <strong>Formular</strong>! 
                    Для завершения регистрации подтвердите свой email-адрес, используя код ниже.
                </p>
                
                <!-- Verification code block -->
                <div style="background: #f8f9ff; border-radius: 12px; padding: 25px; text-align: center; border: 1px dashed #4361ee; margin: 35px 0;">
                    <div style="font-size: 15px; color: #6c757d; margin-bottom: 15px; text-transform: uppercase; letter-spacing: 1px;">Ваш код подтверждения</div>
                    <div style="font-family: 'Courier New', monospace; font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #3a0ca3; word-spacing: 15px;">
                        {{.Code}}
                    </div>
                    <div style="margin-top: 20px; color: #6c757d; font-size: 14px;">
                        Действителен в течение 10 минут
                    </div>
                </div>
                
                <p style="font-size: 16px; margin-bottom: 25px; color: #495057;">
                    Просто введите этот код в приложении Formular, чтобы завершить регистрацию.
                </p>
                
                <p style="font-size: 16px; color: #858585; font-style: italic; border-left: 3px solid #e9ecef; padding-left: 15px; margin: 30px 0;">
                    Если вы не запрашивали подтверждение email, пожалуйста, проигнорируйте это письмо. 
                    Ваш аккаунт останется в безопасности.
                </p>
                
                <!-- Details footer -->
                <div style="background-color: #f8f9ff; border-radius: 10px; padding: 20px; margin-top: 30px; font-size: 14px; color: #6c757d;">
                    <table width="100%">
                        <tr>
                            <td width="50%" style="vertical-align: top;">
                                <strong>Дата отправки:</strong><br>
                                {{.Date}}
                            </td>
                            <td width="50%" style="vertical-align: top; text-align: right;">
                                <strong>Ваш email:</strong><br>
                                {{.User}}
                            </td>
                        </tr>
                    </table>
                </div>
            </td>
        </tr>
        
        <!-- Footer -->
        <tr>
            <td style="background-color: #2b2d42; color: #adb5bd; padding: 30px; text-align: center; font-size: 14px;">
                <div style="margin-bottom: 15px;">
                    <span style="display: inline-block; width: 30px; height: 3px; background: #4361ee; border-radius: 1.5px;"></span>
                </div>
                <p style="margin: 10px 0;">
                    С уважением,<br>
                    <strong style="color: white;">Команда Formular</strong>
                </p>
                <p style="margin: 20px 0 10px; opacity: 0.7; font-size: 12px;">
                    © Formular. Все права защищены.<br>
                    Подготовка к экзаменам с профессиональным подходом
                </p>
                <div style="margin-top: 20px; opacity: 0.5; font-size: 12px;">
                    Это письмо сгенерировано автоматически. Пожалуйста, не отвечайте на него.
                </div>
            </td>
        </tr>
    </table>
    
    <!-- Decorative math elements -->
    <div style="position: absolute; top: 10%; right: 5%; opacity: 0.05; font-size: 48px; transform: rotate(-15deg);">∫</div>
    <div style="position: absolute; bottom: 15%; left: 5%; opacity: 0.05; font-size: 42px; transform: rotate(10deg);">∑</div>
    <div style="position: absolute; bottom: 25%; right: 8%; opacity: 0.03; font-size: 36px;">∞</div>
</body>
</html>
`
