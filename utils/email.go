package utils

import (
	"fmt"
	"net/smtp"
	"time"
)

// إعدادات الإيميل (ضع بياناتك الصحيحة هنا)
const (
	SMTPHost    = "smtp.gmail.com"
	SMTPPort    = "587"
	SenderEmail = "yasserbadr76@gmail.com" // إيميلك
	SenderPass = "vnkbrcvndzqmjlue"    // ⚠️ كلمة السر بدون مسافات نهائيا
)

func SendConfirmationEmail(toEmail, studentName, zoomLink, magicToken string) error {
	// 1. تجهيز ملف التقويم (Calendar Invite)
	meetingTime := time.Now().Add(24 * time.Hour)
	icsContent := fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//QuranApp//V1//EN
BEGIN:VEVENT
UID:%d@quranapp
DTSTAMP:%s
DTSTART:%s
DTEND:%s
SUMMARY:حصة تحفيظ قرآن
DESCRIPTION:رابط الزوم: %s
LOCATION:Zoom
END:VEVENT
END:VCALENDAR`, time.Now().Unix(), time.Now().Format("20060102T150405"), meetingTime.Format("20060102T150405"), meetingTime.Add(1*time.Hour).Format("20060102T150405"), zoomLink)

	// 2. نص الرسالة
	subject := "تأكيد الحجز وتفاصيل الحصة"
	body := fmt.Sprintf(`
السلام عليكم %s،

تم تأكيد حجزك بنجاح!
رابط الزوم: %s

رابط متابعة الدرجات الخاص بك:
https://your-app-url.com/student/%s

مرفق ملف التقويم لإضافته لجدولك.
`, studentName, zoomLink, magicToken)

	// 3. بناء هيكل الإيميل (Header + Body + Attachment)
	boundary := "my-boundary-123"
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/mixed; boundary=%s\r\n\r\n"+
		"--%s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
		"%s\r\n\r\n"+
		"--%s\r\n"+
		"Content-Type: text/calendar; method=REQUEST; name=\"invite.ics\"\r\n"+
		"Content-Transfer-Encoding: base64\r\n"+
		"Content-Disposition: attachment; filename=\"invite.ics\"\r\n\r\n"+
		"%s\r\n"+
		"--%s--", toEmail, subject, boundary, boundary, body, boundary, icsContent, boundary))

	// 4. الإرسال الفعلي
	auth := smtp.PlainAuth("", SenderEmail, SenderPass, SMTPHost)
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SenderEmail, []string{toEmail}, msg)
	
	return err
}