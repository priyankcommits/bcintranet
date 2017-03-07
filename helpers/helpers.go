package helpers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"bcintranet/models"

	"github.com/jung-kurt/gofpdf"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/mgo.v2/bson"
)

func ImageToBase64(url string) string {
	// Convert url image to base64 encoding
	res, _ := http.Get(url)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	imgBase64Str := base64.StdEncoding.EncodeToString(bodyBytes)
	return imgBase64Str
}

func ConvertFormDate(value string) reflect.Value {
	//Converts html date strings to a date type format and returns it
	s, _ := time.Parse("2006-01-02", value)
	return reflect.ValueOf(s.UTC())
}

func SendGridEmail(email *models.Email) error {
	// Send email via send grid
	from := mail.NewEmail("noreply@beautifulcode.in", "intranet@beautifulcode.in")
	subject := email.Subject
	to := mail.NewEmail(email.To.FirstName+email.To.LastName, email.To.Email)
	content := mail.NewContent("text/html", email.Body)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	return err
}

func GeneratePayslipPDF(payslip models.Payslip) {
	// generate PDF for payslip
	salary := payslip.GrossAnnualSalary
	var floatInt int
	floatInt = 2
	var floatType int
	floatType = 64
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetX(-60)
	pdf.SetFont("Arial", "", 16)
	pdf.SetTextColor(26, 162, 251)
	pdf.Cell(30, 0, "BEAUTIFUL ")
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(30, 0, " CODE")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(10)
	pdf.SetXY(100, 30)
	pdf.Line(10, 20, 200, 20)
	pdf.Line(10, 40, 200, 40)
	pdf.Cell(100, 0, "Pay Slip")
	pdf.Line(10, 20, 10, 40)
	pdf.Line(200, 20, 200, 40)
	pdf.SetXY(100, 40)
	pdf.Line(10, 40, 200, 40)
	pdf.Line(10, 70, 200, 70)
	pdf.Line(10, 40, 10, 70)
	pdf.Line(200, 40, 200, 70)
	pdf.SetXY(20, 40)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Pay Period: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Month.Format(time.RFC1123)[7:16])
	pdf.SetXY(100, 40)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Pay Date: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Day.Format(time.RFC1123)[4:16])
	pdf.SetXY(20, 50)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Employee Name: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Name)
	pdf.SetXY(100, 50)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Position: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Position)
	pdf.SetXY(20, 60)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Employee No: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.EmployeeNo)
	pdf.SetXY(100, 70)
	pdf.Line(10, 70, 200, 70)
	pdf.Line(10, 120, 200, 120)
	pdf.Line(10, 70, 10, 120)
	pdf.Line(120, 70, 120, 120)
	pdf.Line(200, 70, 200, 120)
	pdf.SetXY(20, 70)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(70, 10, "Earnings & Allowances")
	pdf.Cell(30, 10, "INR")
	pdf.SetXY(20, 80)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(70, 10, "Basic Salary")
	pdf.Cell(30, 10, strconv.FormatFloat(((salary*0.6)/12), 'f', floatInt, floatType))
	pdf.SetXY(20, 90)
	pdf.Cell(70, 10, "House Rent Allowance")
	pdf.Cell(30, 10, strconv.FormatFloat(((salary*0.2)/12), 'f', floatInt, floatType))
	pdf.SetXY(20, 100)
	pdf.Cell(70, 10, "Spcial / Conv Allowance")
	pdf.Cell(30, 10, strconv.FormatFloat(((salary*0.15)/12), 'f', floatInt, floatType))
	pdf.SetXY(20, 110)
	pdf.Cell(70, 10, "Other Allowance")
	pdf.Cell(30, 10, strconv.FormatFloat(((salary*0.05)/12), 'f', floatInt, floatType))
	pdf.SetXY(120, 70)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Deductions")
	pdf.Cell(20, 10, "INR")
	pdf.SetXY(120, 80)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Income Tax")
	pdf.Cell(20, 10, strconv.FormatFloat(payslip.TDS, 'f', floatInt, floatType))
	pdf.SetXY(120, 90)
	pdf.Cell(40, 10, "Advance")
	pdf.Cell(20, 10, "0.00")
	pdf.SetXY(120, 100)
	pdf.Cell(40, 10, "Profession Tax")
	pdf.Cell(20, 10, "0.00")
	pdf.SetXY(100, 120)
	pdf.Line(10, 120, 200, 120)
	pdf.Line(10, 160, 200, 160)
	pdf.Line(10, 120, 10, 160)
	pdf.Line(120, 120, 120, 160)
	pdf.Line(200, 120, 200, 160)
	pdf.SetXY(20, 120)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(20, 10, "Bank Account: ")
	pdf.SetXY(20, 130)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Account No: ")
	pdf.Cell(50, 10, payslip.AccountNo)
	pdf.SetXY(20, 140)
	pdf.Cell(40, 10, "IFSC Code: ")
	pdf.Cell(50, 10, payslip.IFSCCode)
	pdf.SetXY(120, 120)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Pay Summary")
	pdf.Cell(30, 10, "INR")
	pdf.SetXY(120, 130)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Total Gross")
	pdf.Cell(20, 10, strconv.FormatFloat((payslip.Amount), 'f', floatInt, floatType))
	pdf.SetXY(120, 140)
	pdf.Cell(40, 10, "Deductions")
	pdf.Cell(20, 10, strconv.FormatFloat((payslip.TDS), 'f', floatInt, floatType))
	pdf.SetXY(120, 150)
	pdf.Cell(40, 10, "NET PAY")
	pdf.Cell(20, 10, strconv.FormatFloat((payslip.Amount-payslip.TDS), 'f', floatInt, floatType))
	pdf.SetXY(10, 160)
	pdf.Line(10, 190, 200, 190)
	pdf.Line(10, 160, 10, 190)
	pdf.Line(200, 160, 200, 190)
	pdf.SetXY(75, 170)
	pdf.Cell(150, 10, "(*) denotes back pay adjustment")
	pdf.SetXY(75, 180)
	pdf.Cell(150, 10, "Computer Generated Form does not require signature")
	pdf.OutputFileAndClose("media/" + payslip.PayslipID.Hex() + ".pdf")
}

func InsertAboutData() models.About {
	// Insert Intranet about data
	var about models.About
	about.AboutID = bson.NewObjectId()
	about.Title = "BC Intranet"
	about.About = `BC Intranet is a internal tool for the employees of BC.
	It is implemented to be a centralized source of data that is bound within the company.
	It also useful to be a portal to do daily chores in a time saving manner`
	about.Aim = `I (PK) wanted to learn GoLang and MongoDB by imlementing a real life application`
	about.TechStack = `
	1. GoLang
	2. MongoDB
	3. MaterializeCSS + MaterialCSS
	4. jQuery
	5. Notable Packages: Gorilla Toolkit
	`
	about.CurrentFeatures = `
	1. BC Wall / Announcments
	2. User Profile / Listing
	3. Notifications
	4. PaySlip generator
	5. Project Showcase
	6. Metric Tracking
	`
	about.FutureFeatures = `
	1. Slack Integration
	2. Google Calendar Integration
	3. Direct / Group In App Messaging
	4. More User/ Ajax Interactions whereever possible
	5. More security Implementation(CSRF, XSS, etc)
	`
	about.Limitations = `
	1. Ability to delete / edit certains entities
	2. Email restrictions due to SendGrid
	`
	about.RepoLink = "https://github.com/priyankcommits/bcintranet"
	about.HostLink = "https://fierce-thicket-56115.herokuapp.com/"
	return about
}
