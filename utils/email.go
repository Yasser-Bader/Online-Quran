package utils

import (
	"fmt"
	"net/smtp"
	//"strings"
	"time"
)

// إعدادات الإيميل (يفضل وضعها في .env لاحقاً)
const (
	SMTPHost = "smtp.gmail.com"
	SMTPPort = "587"
	SenderEmail = "YOUR_EMAIL@gmail.com" // ضع إيميلك هنا
	SenderPass  = "YOUR_APP_PASSWORD"    // ضع كلمة مرور التطبيق هنا
)

func SendConfirmationEmail(toEmail, studentName, zoomLink, magicToken string) error {
	// 1. تجهيز ملف التقويم ICS
	meetingTime := time.Now().Add(24 * time.Hour) // مثال: الموعد غداً (يمكن تغييره)
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

	// 2. تجهيز محتوى الإيميل
	subject := "تأكيد الحجز وتفاصيل الحصة"
	body := fmt.Sprintf(`
مرحباً %s،

تم تأكيد حجزك واستلام الدفع بنجاح!
رابط الزوم: %s
رابط متابعة درجاتك (سري): https://your-site.com/student/%s

مرفق ملف التقويم لإضافته لجدولك.
`, studentName, zoomLink, magicToken)

	// 3. بناء الرسالة (MIME) لإرفاق الملف
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
		"--%s--", toEmail, subject, boundary, boundary, body, boundary, icsContent, boundary)) // ملاحظة: في الإنتاج يفضل تشفير icsContent بـ Base64 فعلياً، لكن للتجربة النصية قد يعمل

	// 4. الإرسال
	auth := smtp.PlainAuth("", SenderEmail, SenderPass, SMTPHost)
	return smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SenderEmail, []string{toEmail}, msg)
}