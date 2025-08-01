package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed "Inter/Inter Variable/Inter.ttf"
var interFont []byte

//go:embed "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf"
var interBoldFont []byte

type Invoice struct {
	Id    string `json:"id" yaml:"id"`
	Title string `json:"title" yaml:"title"`

	Logo string `json:"logo" yaml:"logo"`
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
	Date string `json:"date" yaml:"date"`
	Due  string `json:"due" yaml:"due"`

	Services []string  `json:"services" yaml:"services"`
	Rates    []float64 `json:"rates" yaml:"rates"`

	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Note string `json:"note" yaml:"note"`
}

func DefaultInvoice() Invoice {
	return Invoice{
		Id:       time.Now().Format("20060102"),
		Title:    "INVOICE",
		Rates:    []float64{25},
		Services: []string{"Software Development"},
		From:     "Project Folded, Inc.",
		To:       "Untitled Corporation, Inc.",
		Date:     time.Now().Format("Jan 02, 2006"),
		Due:      time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Discount: 0,
		Currency: "USD",
	}
}

var (
	importPath     string
	output         string
	file           = Invoice{}
	defaultInvoice = DefaultInvoice()
)

func init() {
	viper.AutomaticEnv()

	generateCmd.Flags().StringVar(&importPath, "import", "", "Imported file (.json/.yaml)")
	generateCmd.Flags().StringVar(&file.Id, "id", time.Now().Format("20060102"), "ID")
	generateCmd.Flags().StringVar(&file.Title, "title", "INVOICE", "Title")

	generateCmd.Flags().Float64SliceVarP(&file.Rates, "rate", "r", defaultInvoice.Rates, "Rates")
	generateCmd.Flags().StringSliceVarP(&file.Services, "service", "s", defaultInvoice.Services, "Services")

	generateCmd.Flags().StringVarP(&file.Logo, "logo", "l", defaultInvoice.Logo, "Company logo")
	generateCmd.Flags().StringVarP(&file.From, "from", "f", defaultInvoice.From, "Issuing company")
	generateCmd.Flags().StringVarP(&file.To, "to", "t", defaultInvoice.To, "Recipient company")
	generateCmd.Flags().StringVar(&file.Date, "date", defaultInvoice.Date, "Date")
	generateCmd.Flags().StringVar(&file.Due, "due", defaultInvoice.Due, "Payment due date")

	generateCmd.Flags().Float64VarP(&file.Discount, "discount", "d", defaultInvoice.Discount, "Discount")
	generateCmd.Flags().StringVarP(&file.Currency, "currency", "c", defaultInvoice.Currency, "Currency")

	generateCmd.Flags().StringVarP(&file.Note, "note", "n", "", "Note")
	generateCmd.Flags().StringVarP(&output, "output", "o", "invoice.pdf", "Output file (.pdf)")

	flag.Parse()
}

var rootCmd = &cobra.Command{
	Use:   "invoice",
	Short: "Invoice generates invoices from the command line.",
	Long:  `Invoice generates invoices from the command line.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an invoice",
	Long:  `Generate an invoice`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if importPath != "" {
			err := importData(importPath, &file, cmd.Flags())
			if err != nil {
				return err
			}
		}

		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{
			PageSize: *gopdf.PageSizeA4,
		})
		pdf.SetMargins(40, 40, 40, 40)
		pdf.AddPage()
		err := pdf.AddTTFFontData("Inter", interFont)
		if err != nil {
			return err
		}

		err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
		if err != nil {
			return err
		}

		writeLogo(&pdf, file.Logo, file.From)
		writeTitle(&pdf, file.Title, file.Id, file.Date)
		writeBillTo(&pdf, file.To)
		writeHeaderRow(&pdf)
		subtotal := 0.0
		for i := range file.Services {
			r := 0.0
			if len(file.Rates) > i {
				r = file.Rates[i]
			}

			writeRow(&pdf, file.Services[i], r)
			subtotal += r
		}
		if file.Note != "" {
			writeNotes(&pdf, file.Note)
		}
		writeTotals(&pdf, subtotal, subtotal*file.Discount)
		if file.Due != "" {
			writeDueDate(&pdf, file.Due)
		}
		writeFooter(&pdf, file.Id)
		output = strings.TrimSuffix(output, ".pdf") + ".pdf"
		err = pdf.WritePdf(output)
		if err != nil {
			return err
		}

		fmt.Printf("Generated %s\n", output)

		return nil
	},
}

func main() {
	rootCmd.AddCommand(generateCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
