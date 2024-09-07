package main

func createUrl(url *Url) error {

	_, err := db.Exec("INSERT INTO public.url (original, shorten) VALUES ($1, $2)",
		url.Original,
		url.Shorten,
	)

	return err
}

func getOriginalUrl(shorten string) (string, error) {
	u := ""
	row := db.QueryRow(
		"SELECT original FROM url WHERE shorten=$1;",
		shorten)

	err := row.Scan(&u)

	if err != nil {
		return "", err
	}

	return u, nil
}
