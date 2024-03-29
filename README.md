<div align="center">
<h1>
  <code>notionion</code> 🧅
</h1>
  <img src="https://github.com/ariary/notionion/blob/main/img/onion-logo.png"  width=150>
  
  <strong> Use <a href="https://www.notion.so">Notion</a> as an HTTP proxy.</strong><br>
  <i>🫀Like having Burp in your notetaking app</i>
</div>

---

![demo](https://github.com/ariary/Notionion/blob/main/img/demo-fast.gif)

---

<div align=left>
<h3 >Why? 🤔 </h3>
Mainly for fun of adding a blade to the Swiss army knife that is notion. <br>The purpose is to provide an HTTP proxy which <b>takes advantage of the Notion benefits</b>
<ul>
<li>Cross-platform (Windows/MacOS, browsers, mobile)</li>
<li>Shared with authentication</li>
<li>Accomplished notetaking app (GUI provided, ease HTTP packet manipulation, add notes)</li>
</ul>
</div>
<div align=right>
<h3 >How?  🤷‍♂️</h3>
Just use notion as usual and launch <code>notionion</code>.
</div>

---
<blockquote align=left>
roughly inspired by the great idea of <a href="https://github.com/mttaggart/OffensiveNotion">OffensiveNotion</a>! 
</blockquote>

## Quickstart

**Set-up**  ([details](#-set-up))
1. Create the "Proxy page" in Notion
2. Give the permissions to `notionion` to access the Proxy page 

**Run** ([details](#-run))

3. Perform HTTP request
4. Modify it trough the "Proxy page" in Notion
5. See result of request

### 🏗️ Set-up

#### Create the "Proxy page" in Notion

<img src="https://github.com/ariary/Notionion/blob/main/img/proxy-page.png"  width=500>

You can duplicate the template [notionion template](https://fluff-grade-468.notion.site/notionion_template-f95213ec89a04f66ad895ddac850d33e)

#### Give the permissions to `notionion` to access the Proxy page
* Go to the [Notion API developer page](https://developers.notion.com/) and log in. Create an Integration user (`New integration`). Copy that user's API key
* Copy the "Proxy page" Url
  * In browser: only copy the URL
  * On desktop app: `CTRL+L`
* Add your Notion Developer API account to this page (In the upper-right corner of your Notion page, click ***"Share"*** and ***"Invite"***)
* Install `notionion` [see](#install)

#### Declare environment variables to specify the notion proxy page:
```shell
source env.sh
# Alternatively, you can just export NOTION_TOKEN (which is the api key) & NOTION_PAGE_URL
```

### 👟 Run

```shell
notionion
```


## Install
* **From release**: `curl -lO -L https://github.com/ariary/notionion/releases/latest/download/notionion && chmod +x notionion`
* **Build it**: `git clone https://github.com/ariary/notionion && make before.build && make build.notion` *(need `go`)*
* **with `go`**:`go install github.com/ariary/notionion@latest`
