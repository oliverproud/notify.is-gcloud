package statements

const (
	// InstagramUpdateStatement updates the instagram column and timestamp if a username is available
	InstagramUpdateStatement = `
  UPDATE users
  SET instagram = false, timestamp = (now() at time zone 'utc')
  WHERE id = $1;
  `
	// TwitterUpdateStatement updates the twitter column and timestamp if a username is available
	TwitterUpdateStatement = `
  UPDATE users
  SET twitter = false, timestamp = (now() at time zone 'utc')
  WHERE id = $1;
  `
	// GithubUpdateStatement updates the github column and timestamp if a username is available
	GithubUpdateStatement = `
  UPDATE users
  SET github = false, timestamp = (now() at time zone 'utc')
  WHERE id = $1;
  `
	// DefaultUpdateStatement updates the github column and timestamp if a username is available
	DefaultUpdateStatement = `
  UPDATE users
  SET timestamp = (now() at time zone 'utc')
  WHERE id = $1;
  `

	// SelectStatement retrieves records that haven't been checked for over 12 hours
	SelectStatement = `
	SELECT id, first_name, email, username, instagram, twitter, github, timestamp
	FROM users
	WHERE EXTRACT(EPOCH FROM ((now() at time zone 'utc') - timestamp)) > 43000.0;
	`
)
