package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	productDomain "erp-api/internal/domain/product"
	productUseCase "erp-api/internal/usecase/product"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// Handler lida com exportação de relatórios usando os DTOs existentes de produto
type Handler struct {
	productUseCase productUseCase.UseCaseInterface
}

// NewHandler cria um novo handler de reports
func NewHandler(productUseCase productUseCase.UseCaseInterface) *Handler {
	return &Handler{productUseCase: productUseCase}
}

// Export exporta lista de produtos em pdf, xlsx ou preview json
// @Summary Exportar relatório de produtos
// @Description Exporta produtos em PDF, Excel ou preview JSON usando DTOs existentes
// @Tags reports
// @Produce json,application/pdf,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param format query string true "Formato" Enums(pdf,xlsx,preview)
// @Param limit query int false "Limite (default 100)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} productDomain.ProductListDTO
// @Router /api/reports/export [get]
func (h *Handler) Export(c *gin.Context) {
	format := c.Query("format")
	if format == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format parameter is required"})
		return
	}

	limit := parseIntDefault(c.Query("limit"), 100)
	offset := parseIntDefault(c.Query("offset"), 0)

	listDTO, err := h.buildProductListDTO(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load products"})
		return
	}

	switch format {
	case "pdf":
		h.exportPDF(c, listDTO)
	case "xlsx":
		h.exportExcel(c, listDTO)
	case "preview":
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, listDTO)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format. Use pdf, xlsx or preview"})
	}
}

func (h *Handler) buildProductListDTO(c *gin.Context, limit, offset int) (*productDomain.ProductListDTO, error) {
	products, err := h.productUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := h.productUseCase.Count(c.Request.Context())
	if err != nil {
		return nil, err
	}

	resp := &productDomain.ProductListDTO{
		Products: make([]*productDomain.ProductDTO, len(products)),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}
	for i, p := range products {
		resp.Products[i] = &productDomain.ProductDTO{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Type:        p.Type,
			Price:       p.Price,
			Unit:        p.Unit,
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		}
	}
	return resp, nil
}

func (h *Handler) exportPDF(c *gin.Context, data *productDomain.ProductListDTO) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(60, 10, "Relatório de Produtos")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 6, fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04:05")))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "Nome", "Tipo", "Unid", "Preço", "Ativo", "Criado em", "Atualizado em"}
	widths := []float64{25, 70, 25, 18, 25, 15, 40, 40}
	for i, hText := range headers {
		pdf.CellFormat(widths[i], 7, hText, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, p := range data.Products {
		row := []string{
			p.ID,
			p.Name,
			p.Type,
			p.Unit,
			fmt.Sprintf("R$ %.2f", p.Price),
			boolToStr(p.IsActive),
			p.CreatedAt,
			p.UpdatedAt,
		}
		for i, v := range row {
			pdf.CellFormat(widths[i], 6, truncate(v, 45), "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	pdf.SetFont("Arial", "B", 10)
	pdf.Ln(4)
	pdf.Cell(40, 6, fmt.Sprintf("Total de registros: %d", data.Total))

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=produtos.pdf")
	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PDF"})
	}
}

func (h *Handler) exportExcel(c *gin.Context, data *productDomain.ProductListDTO) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Produtos"
	_, _ = f.NewSheet(sheet)
	f.SetActiveSheet(1)

	f.SetCellValue(sheet, "A1", "Relatório de Produtos")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04:05")))

	headers := []string{"ID", "Nome", "Descrição", "Tipo", "Unidade", "Preço", "Ativo", "Criado em", "Atualizado em"}
	for i, hText := range headers {
		col, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, col, hText)
	}

	for i, p := range data.Products {
		row := i + 5
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Description)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.Type)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), p.Unit)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.Price)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.IsActive)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), p.CreatedAt)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), p.UpdatedAt)
	}

	totalRow := len(data.Products) + 6
	f.SetCellValue(sheet, fmt.Sprintf("H%d", totalRow), "Total Registros:")
	f.SetCellValue(sheet, fmt.Sprintf("I%d", totalRow), data.Total)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=produtos.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate Excel"})
	}
}

// Helpers
func parseIntDefault(v string, def int) int {
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}

func boolToStr(b bool) string {
	if b {
		return "Sim"
	}
	return "Não"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
