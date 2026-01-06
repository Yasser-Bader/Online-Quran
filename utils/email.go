package utils

import (
	"fmt"
	"net/smtp"
)

// إعدادات الإيميل
const (
	SMTPHost = "smtp.gmail.com"
	SMTPPort = "587"
	// تم تعديل الإيميل ليتوافق مع كلمة المرور (App Password) الخاصة بـ yasserbadr76
	SenderEmail = "yasserbadr76@gmail.com"
	SenderPass  = "vnkbrcvndzqmjlue" // ⚠️ كلمة السر بدون مسافات نهائيا

)

/*
func SendConfirmationEmail(toEmail, studentName, zoomLink, magicToken string) error {
	// 1. تجهيز ملف التقويم (Calendar Invite)
	meetingTime := time.Now().Add(24 * time.Hour)
	timeFormat := "20060102T150405" // تنسيق الوقت المطلوب لملفات iCalendar

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
END:VCALENDAR`,
		time.Now().Unix(),
		time.Now().Format(timeFormat),
		meetingTime.Format(timeFormat),
		meetingTime.Add(1*time.Hour).Format(timeFormat),
		zoomLink)

	// --- تعديل 1: تشفير محتوى ملف التقويم ---
	// هذا ضروري جداً لكي لا يصل الملف تالفاً
	encodedICS := base64.StdEncoding.EncodeToString([]byte(icsContent))

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

	// --- تعديل 2: إضافة From Header ---
	// عدم وجود From هو السبب الرئيسي لدخول الرسائل في الـ Spam
	header := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/mixed; boundary=%s\r\n\r\n", SenderEmail, toEmail, subject, boundary)

	bodyPart := fmt.Sprintf("--%s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
		"%s\r\n\r\n", boundary, body)

	// استخدام المحتوى المشفر encodedICS هنا
	attachmentPart := fmt.Sprintf("--%s\r\n"+
		"Content-Type: text/calendar; method=REQUEST; name=\"invite.ics\"\r\n"+
		"Content-Transfer-Encoding: base64\r\n"+
		"Content-Disposition: attachment; filename=\"invite.ics\"\r\n\r\n"+
		"%s\r\n"+
		"--%s--", boundary, encodedICS, boundary)

	msg := []byte(header + bodyPart + attachmentPart)

	// 4. الإرسال الفعلي
	auth := smtp.PlainAuth("", SenderEmail, SenderPass, SMTPHost)
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SenderEmail, []string{toEmail}, msg)

	return err
}*/

// ... (نفس الثوابت السابقة) ...

// 1. دالة الرفض (جديدة)
func SendRejectionEmail(toEmail, name, reason string) error {
	subject := "تحديث بخصوص طلب التسجيل"
	body := fmt.Sprintf(`
مرحباً %s،

نأسف لإبلاغك بأنه تم إلغاء طلب حجزك وحذف بياناتك من النظام.
السبب: %s

يمكنك إعادة التسجيل مرة أخرى ببيانات صحيحة.
`, name, reason)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", SenderEmail, SenderPass, SMTPHost)
	return smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SenderEmail, []string{toEmail}, msg)
}

// 2. دالة القبول (معدلة لتعرف هل هو زوم أم واتساب)
func SendConfirmationEmail(toEmail, studentName, link, magicToken, mode string) error {
	linkType := "رابط الزوم"
	if mode == "group" {
		linkType = "رابط مجموعة الواتساب"
	}

	subject := "تم تأكيد الحجز بنجاح"
	body := fmt.Sprintf(`
مرحباً %s،

تم تأكيد حجزك بنجاح!
%s الخاص بك: %s

رابط متابعة الدرجات: https://your-app.com/student/%s
`, studentName, linkType, link, magicToken)

	// ... (باقي كود الإرسال كما هو) ...
	// اختصاراً للكود هنا، استخدم نفس طريقة الإرسال السابقة
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", SenderEmail, SenderPass, SMTPHost)
	return smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SenderEmail, []string{toEmail}, msg)
}
