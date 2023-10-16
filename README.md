# alist-northflank-postgresql


## Deploy alist to northflank
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/dreamoeu/alist-northflank-postgres)

Use northflank's add-on postgres database, your settings will be persistent, don't worry about hibernate losing configuration.

If you can't deploy, and it says "We couldn't deploy your app because the source code violates the Salesforce Acceptable Use and External-Facing Services Policy.", you need to fork this repo and click the `Deploy` button in your own fork.

## Variables

Here are explanations for the variables set during deployment:

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `CDN` | | CDN address. If you want to use a CDN, you can set this field. `$version` will be replaced with the actual version of `alist-web`. |
| `DATABASE_URL` | | Database connection URL. It defaults to an empty value, indicating the use of a local SQLite database. |
| `SITE_URL` | | Website URL, for example, https://example.com. This address is used in certain parts of the Alist program, and if not set, some features may not work correctly. |
| `TZ` | `Asia/Shanghai` | Timezone, with the China timezone being Asia/Shanghai. |

For more variables and explanations, please refer to:

https://github.com/alist-org/alist/blob/main/internal/conf/config.go

https://alist.nn.ci/config/configuration.html

## Get Password
`More` -> `View logs` -> You will see your password, if it is scrolled to the top and out of view, click `Restart all dynos` then the log will be redisplayed.
