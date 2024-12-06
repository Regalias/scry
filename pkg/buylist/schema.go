package buylist

const buylistTableName = "buylists"
const cardsTableName = "cards"
const selectionsTableName = "selections"

const buylistSchema = `
CREATE TABLE IF NOT EXISTS buylists (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at NUMERIC NOT NULL
);
`

const cardSchema = `
CREATE TABLE IF NOT EXISTS cards (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	quantity INTEGER NOT NULL,
	buylist_id INTEGER NOT NULL,
	CONSTRAINT fk_buylists
		FOREIGN KEY (buylist_id)
		REFERENCES buylists(id)
		ON DELETE CASCADE
);
`

const selectionsSchema = `
CREATE TABLE IF NOT EXISTS selections (
	id INTEGER PRIMARY KEY,
	quantity INTEGER NOT NULL,
	is_purchased INTEGER NOT NULL,
	is_flagged INTEGER NOT NULL,
	offering BLOB NOT NULL,
	card_id INTEGER NOT NULL,
	CONSTRAINT fk_cards
		FOREIGN KEY (card_id)
		REFERENCES cards(id)
		ON DELETE CASCADE
);
`
