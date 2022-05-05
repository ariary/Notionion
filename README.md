<div align="center">
<h1>
  <code>notionion</code>🧅
</h1>
  <img src="https://github.com/ariary/notionion/blob/main/img/onion-logo.png"  width=150>
  
  <strong> Use <a href="https://www.notion.so">Notion</a> as an HTTP proxy.</strong><br>
  <i>🫀Like having Burp in your notetaking app</i>
</div>

<blockquote align=left>
roughly inspired by the great idea of <a href="https://github.com/mttaggart/OffensiveNotion">OffensiveNotion</a>! 
</blockquote>

---

... DEMO

---

## Quickstart

**Set-up** ([details]())
1. Create the "Interception page" in Notion
2. Give the permissions to `notionion` to access the Interception page 

**Run** ([details]())

3. Perform HTTP request
4. Modify it trough the "Interception page" in Notion
5. See result of request

### 🏗️ Set-up

#### Create the "Interception page" in Notion

<img src="https://github.com/ariary/notionion/blob/main/img/page.png"  width=350>

<sup><i>You can import the template [`./page.html`](https://github.com/ariary/notionion/blob/main/page.html) using the "Import" function of Notion<i></sup>

#### Give the permissions to `notionion` to access the Interception page
* Go to the [Notion API developer page](https://developers.notion.com/) and log in. Create an Integration user (`New integration`). Copy that user's API key
* Copy the "Interception page" Url
  * In browser: only copy the URL
  * On desktop app: `CTRL+L`
* Add your Notion Developer API account to this page (In the upper-right corner of your Notion page, click ***"Share"*** and ***"Invite"***)
* Install `notionion`:
  * `git clone https://github.com/ariary/notionion && make before.build && make build.notion` *(need `go`)*
  * `go install github.com/ariary/notionion@latest`
  * `curl -lO -L https://github.com/ariary/notionion/releases/latest/download/notionion && chmod +x notionion`

### 👟 Run
