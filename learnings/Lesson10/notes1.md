# Setting Security Headers

1. [Content-Security-Policy (often abbreviated to CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)

2. [Referrer Policy](https://developer.mozilla.org/en-US/docs/Web/Security/Same-origin_policy): 

### What is the `Referrer-Policy`?

When you click a link or make a request from one webpage to another, the browser often sends the URL of the page you're coming from. This is called the "referrer" (note: "referer" is a common misspelling that has stuck around in web terminology). The `Referrer-Policy` is a setting that determines how much of this referrer information should be shared.

### Different `Referrer-Policy` Values:

Here are some key values you can use for the `Referrer-Policy`:

1. **`no-referrer`**:
   - No referrer information is sent with any requests.

2. **`no-referrer-when-downgrade`**:
   - Referrer information is sent for same-origin and secure (HTTPS) requests, but not for secure (HTTPS) to insecure (HTTP) requests.

3. **`origin`**:
   - Only the origin (scheme, host, and port) of the referrer URL is sent, not the full path. For example, if the referrer URL is `https://example.com/page`, only `https://example.com` will be sent.

4. **`origin-when-cross-origin`**:
   - The full URL is sent as the referrer for same-origin requests, but only the origin is sent for cross-origin requests.

5. **`same-origin`**:
   - The full referrer URL is sent for same-origin requests, but no referrer information is sent for cross-origin requests.

6. **`strict-origin`**:
   - Only the origin is sent as the referrer, but only when navigating from HTTPS to HTTPS. No referrer information is sent otherwise.

7. **`strict-origin-when-cross-origin`**:
   - The full URL is sent for same-origin requests and only the origin for cross-origin requests, but only when both URLs are HTTPS. If navigating from HTTPS to HTTP, no referrer information is sent.

8. **`unsafe-url`**:
   - The full URL is always sent as the referrer, regardless of the request context.

### What Does `origin-when-cross-origin` Mean?

- **Same-origin requests:** When you navigate within the same website, the full URL of the referrer is sent.
  - Example: Navigating from `https://example.com/page1` to `https://example.com/page2` sends `https://example.com/page1` as the referrer.

- **Cross-origin requests:** When you navigate to a different website, only the origin (scheme, host, and port) of the referrer is sent.
  - Example: Navigating from `https://example.com/page1` to `https://otherwebsite.com/page` sends `https://example.com` as the referrer.

### Why Use `origin-when-cross-origin`?

1. **Privacy:** When you visit a different website, it won't get the full URL of the page you came from, only the origin. This limits the amount of information shared with third-party sites.
2. **Functionality:** Within your own website, you still get the full URL, which can be useful for analytics, debugging, and other purposes.

##

> Again reminding, `http.Handler` in literal sense means that, the one who handles/serves(सेवा करना) http requests, when any request arrives, it serves its requirements. In Go terms, `http.Handler` is an interface which has method `ServeHTTP`
```
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```