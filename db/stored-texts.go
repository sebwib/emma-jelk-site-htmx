package db

import (
	"time"

	"github.com/google/uuid"
)

type StoredText struct {
	UUID        string
	ReferenceID string
	Content     string
	CreatedAt   string
}

func (db *DB) createStoredTextsTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS stored_texts (
		uuid TEXT PRIMARY KEY,
		reference_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TEXT NOT NULL
	);
	`)

	if err != nil {
		return err
	}

	return db.ensureDefaultStoredTexts()
}

func (db *DB) ensureDefaultStoredTexts() error {
	defaults := []StoredText{
		{
			ReferenceID: "home_title",
			Content:     "Hem",
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			ReferenceID: "gallery_title",
			Content:     "Galleri",
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			ReferenceID: "buy_art_title",
			Content:     "Köp konst",
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			ReferenceID: "buy_art_text",
			Content:     "Under fliken galleri hittar du mina originalmålningar, både tillgängliga och sålda verk. Är du intresserad av att köpa originalkonst och vill få mer information om ett specifikt verk, skicka ett mail till: emma.jelk@gmail.com",
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			ReferenceID: "about_me_title",
			Content:     "Om mig",
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			ReferenceID: "about_me_text",
			Content: `Under min uppväxt tecknade jag dagligen, och intresset var ett brinnande sådant. Efter barndomen gick intresset för tecknandet i vågor tills jag tillslut fick upp ögonen för tatuering, och därmed hamnade jag som lärling på en lokal tatueringsstudio. Tatuerandet var enormt utvecklande då det ingick i min dagliga arbetsrutin att vara kreativ, uppleva kundkontakt samt arbeta disciplinerat och väldigt noggrant under alla moment.
Det är inte omöjligt att det var den disciplinen som fick mig att så småningom börja måla realism, då det kräver en stor noggrannhet och koncentration för att uppnå det resultat jag vill ha.


Jag arbetade sedan på tatueringsstudion i ett par år, och plötsligt stod familjelivet framför dörren. Återigen hamnade tecknandet lite mer i skymundan mellan barn, arbete, plugg och ännu ett barn. 


2018 bestämde jag mig för att måla en tavla. Jag hade då inte målat med pensel och färg på många år utan istället tecknat/målat digitalt när små stunder av tid fanns. Jag var höggravid, och la över trettio långa timmar på min första “riktiga” målning. När jag tillslut tog några steg tillbaka och betraktade min målning ropade jag lyckligt till min man Sebastian: “Jag kan visst måla!”


Efter den stunden smög sig känslan på mig. Att det var precis det här jag ville göra. Och efter att ha lekt med tanken, känt och klämt på den under ett års tid medan jag var föräldraledig med mitt andra barn, vågade jag hösten 2019 uttala orden: “jag vill bli konstnär.”


Så nu gör jag det här.


Målningarna du ser här är resultatet av mina första ca 200 timmar av akrylmåleri. Början på min konstnärsresa, utforskande av både form, färg och teknik. Den enda röda tråden jag själv lyckats se i mitt målande är att jag vill måla allt. Kanske för att jag tycker mig se något i ett ögonblick, ofta ett väldigt vardagligt sådant, något enormt stort i det vardagliga. Och att få plocka ut just det där ögonblicket, och studera dess detaljer är för mig fantastiskt. Det hjälper mig att se världen på ett sätt som ger mig välbefinnande och närhet till livet självt.`,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	for _, def := range defaults {
		var count int
		err := db.QueryRow(`
		SELECT COUNT(*)
		FROM stored_texts
		WHERE reference_id = ?;
		`, def.ReferenceID).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			if err := db.AddStoredText(def); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *DB) AddStoredText(text StoredText) error {
	text.UUID = uuid.NewString()
	text.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Exec(`
	INSERT INTO stored_texts (uuid, reference_id, content, created_at)
	VALUES (?, ?, ?, ?);
	`, text.UUID, text.ReferenceID, text.Content, text.CreatedAt)
	return err
}

func (db *DB) GetReferences() ([]StoredText, error) {
	rows, err := db.Query(`
	SELECT uuid, reference_id, content, created_at
	FROM stored_texts
	ORDER BY created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []StoredText
	for rows.Next() {
		var reference StoredText
		if err := rows.Scan(&reference.UUID, &reference.ReferenceID, &reference.Content, &reference.CreatedAt); err != nil {
			return nil, err
		}
		references = append(references, reference)
	}
	return references, nil
}

func (db *DB) GetStoredTextByReferenceID(referenceID string) ([]StoredText, error) {
	rows, err := db.Query(`
	SELECT uuid, reference_id, content, created_at
	FROM stored_texts
	WHERE reference_id = ?
	ORDER BY created_at DESC;
	`, referenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var texts []StoredText
	for rows.Next() {
		var text StoredText
		if err := rows.Scan(&text.UUID, &text.ReferenceID, &text.Content, &text.CreatedAt); err != nil {
			return nil, err
		}
		texts = append(texts, text)
	}
	return texts, nil
}
