package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"flashcards/api/internal/domain"
	"flashcards/api/internal/middleware"
	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cards *service.CardService
}

func NewCardHandler(cards *service.CardService) *CardHandler {
	return &CardHandler{cards: cards}
}

func (h *CardHandler) PublicStats(c *gin.Context) {
	n, err := h.cards.CountAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count cards"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_cards": n})
}

type cardBody struct {
	Front     string   `json:"front" binding:"required"`
	Back      string   `json:"back" binding:"required"`
	Topic     *string  `json:"topic"`
	Level     *string  `json:"level"`
	SourceURL *string  `json:"sourceUrl"`
	Tags      []string `json:"tags"`
}

func (b cardBody) toInput() service.CardInput {
	return service.CardInput{
		Front:     b.Front,
		Back:      b.Back,
		Topic:     b.Topic,
		Level:     b.Level,
		SourceURL: b.SourceURL,
		Tags:      b.Tags,
	}
}

func (h *CardHandler) Create(c *gin.Context) {
	var body cardBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := h.cards.Create(c.GetString(middleware.ContextUserID), c.Param("id"), body.toInput())
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, card)
}

func (h *CardHandler) Update(c *gin.Context) {
	var body cardBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := h.cards.Update(c.GetString(middleware.ContextUserID), c.Param("id"), body.toInput())
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) Delete(c *gin.Context) {
	err := h.cards.Delete(c.GetString(middleware.ContextUserID), c.Param("id"))
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type importBody struct {
	Format string `json:"format" binding:"required,oneof=json csv"`
	Data   string `json:"data"   binding:"required"`
}

func (h *CardHandler) Import(c *gin.Context) {
	var body importBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputs, err := parseImport(body.Format, body.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(inputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no cards found in import data"})
		return
	}

	cards, err := h.cards.BulkCreate(
		c.GetString(middleware.ContextUserID),
		c.Param("id"),
		inputs,
	)
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"imported": len(cards), "cards": cards})
}

func parseImport(format, data string) ([]service.CardInput, error) {
	switch format {
	case "json":
		return parseJSONImport(data)
	case "csv":
		return parseCSVImport(data)
	default:
		return nil, errors.New("unsupported format")
	}
}

func parseJSONImport(data string) ([]service.CardInput, error) {
	var rows []struct {
		Front     string   `json:"front"`
		Back      string   `json:"back"`
		Topic     *string  `json:"topic"`
		Level     *string  `json:"level"`
		SourceURL *string  `json:"sourceUrl"`
		Tags      []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(data), &rows); err != nil {
		return nil, err
	}
	inputs := make([]service.CardInput, 0, len(rows))
	for _, r := range rows {
		if strings.TrimSpace(r.Front) == "" || strings.TrimSpace(r.Back) == "" {
			continue
		}
		inputs = append(inputs, service.CardInput{
			Front:     r.Front,
			Back:      r.Back,
			Topic:     r.Topic,
			Level:     r.Level,
			SourceURL: r.SourceURL,
			Tags:      r.Tags,
		})
	}
	return inputs, nil
}

func parseCSVImport(data string) ([]service.CardInput, error) {
	r := csv.NewReader(strings.NewReader(data))
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, nil
	}

	// Build column index from header row.
	colIdx := map[string]int{}
	for i, h := range records[0] {
		colIdx[strings.ToLower(strings.TrimSpace(h))] = i
	}
	frontCol, hasFront := colIdx["front"]
	backCol, hasBack := colIdx["back"]
	if !hasFront || !hasBack {
		return nil, errors.New("CSV must have 'front' and 'back' columns")
	}

	col := func(row []string, name string) *string {
		idx, ok := colIdx[name]
		if !ok || idx >= len(row) {
			return nil
		}
		v := strings.TrimSpace(row[idx])
		if v == "" {
			return nil
		}
		return &v
	}

	inputs := make([]service.CardInput, 0, len(records)-1)
	for _, row := range records[1:] {
		if frontCol >= len(row) || backCol >= len(row) {
			continue
		}
		front := strings.TrimSpace(row[frontCol])
		back := strings.TrimSpace(row[backCol])
		if front == "" || back == "" {
			continue
		}
		in := service.CardInput{
			Front:     front,
			Back:      back,
			Topic:     col(row, "topic"),
			Level:     col(row, "level"),
			SourceURL: col(row, "sourceurl"),
		}
		if idx, ok := colIdx["tags"]; ok && idx < len(row) {
			raw := strings.TrimSpace(row[idx])
			if raw != "" {
				parts := strings.Split(raw, "|")
				tags := make([]string, 0, len(parts))
				for _, p := range parts {
					if t := strings.TrimSpace(p); t != "" {
						tags = append(tags, t)
					}
				}
				in.Tags = tags
			}
		}
		inputs = append(inputs, in)
	}
	return inputs, nil
}

func cardErrorStatus(err error) (int, bool) {
	switch {
	case err == nil:
		return 0, false
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, true
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, true
	default:
		return http.StatusInternalServerError, true
	}
}

func (h *CardHandler) ExportDeck(c *gin.Context) {
	format := c.Query("format")
	contentType, ext, err := exportMeta(format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	de, err := h.cards.ExportDeck(c.GetString(middleware.ContextUserID), c.Param("id"))
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	data, err := renderExport(format, []service.DeckExport{*de})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := "flashcards-" + de.Deck.Title + "." + ext
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

func (h *CardHandler) ExportAll(c *gin.Context) {
	format := c.Query("format")
	contentType, ext, err := exportMeta(format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decks, err := h.cards.ExportAll(c.GetString(middleware.ContextUserID))
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	data, err := renderExport(format, decks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=flashcards-export."+ext)
	c.Data(http.StatusOK, contentType, data)
}

// exportMeta returns content-type and file extension for a given format.
func exportMeta(format string) (contentType, ext string, err error) {
	switch format {
	case "csv":
		return "text/csv", "csv", nil
	case "json":
		return "application/json", "json", nil
	case "anki-csv":
		return "text/plain", "txt", nil
	default:
		return "", "", errors.New("unsupported format: use csv, json, or anki-csv")
	}
}

// renderExport serializes DeckExport slices to bytes in the requested format.
func renderExport(format string, decks []service.DeckExport) ([]byte, error) {
	switch format {
	case "csv":
		return renderCSVExport(decks)
	case "json":
		return renderJSONExport(decks)
	case "anki-csv":
		return renderAnkiCSVExport(decks)
	default:
		return nil, errors.New("unsupported format")
	}
}

func renderCSVExport(decks []service.DeckExport) ([]byte, error) {
	var buf strings.Builder
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"deck_name", "front", "back", "topic", "level", "source_url", "tags"})
	for _, de := range decks {
		for _, c := range de.Cards {
			row := []string{
				de.Deck.Title,
				c.Content.Front,
				c.Content.Back,
				derefStr(c.Content.Topic),
				derefStr(c.Content.Level),
				derefStr(c.Content.SourceURL),
				strings.Join(c.Content.Tags, "|"),
			}
			_ = w.Write(row)
		}
	}
	w.Flush()
	return []byte(buf.String()), w.Error()
}

func renderJSONExport(decks []service.DeckExport) ([]byte, error) {
	type cardRow struct {
		Front     string   `json:"front"`
		Back      string   `json:"back"`
		Topic     *string  `json:"topic,omitempty"`
		Level     *string  `json:"level,omitempty"`
		SourceURL *string  `json:"sourceUrl,omitempty"`
		Tags      []string `json:"tags"`
	}
	type deckRow struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Cards       []cardRow `json:"cards"`
	}
	rows := make([]deckRow, 0, len(decks))
	for _, de := range decks {
		cards := make([]cardRow, 0, len(de.Cards))
		for _, c := range de.Cards {
			cards = append(cards, cardRow{
				Front:     c.Content.Front,
				Back:      c.Content.Back,
				Topic:     c.Content.Topic,
				Level:     c.Content.Level,
				SourceURL: c.Content.SourceURL,
				Tags:      c.Content.Tags,
			})
		}
		rows = append(rows, deckRow{
			Name:        de.Deck.Title,
			Description: de.Deck.Description,
			Cards:       cards,
		})
	}
	return json.MarshalIndent(rows, "", "  ")
}

func renderAnkiCSVExport(decks []service.DeckExport) ([]byte, error) {
	// Anki tab-separated: front TAB back TAB tags (space-separated), no header.
	var buf strings.Builder
	for _, de := range decks {
		for _, c := range de.Cards {
			buf.WriteString(c.Content.Front)
			buf.WriteByte('\t')
			buf.WriteString(c.Content.Back)
			buf.WriteByte('\t')
			buf.WriteString(strings.Join(c.Content.Tags, " "))
			buf.WriteByte('\n')
		}
	}
	return []byte(buf.String()), nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
