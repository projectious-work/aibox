(function () {
  "use strict";

  var menu = document.querySelector(".release-menu__list[data-release-manifest]");
  if (!menu) return;

  function activeBase(manifest) {
    var candidates = [];
    manifest.lines.forEach(function (line) {
      if (!line || !line.current || typeof line.current.url !== "string" || !Array.isArray(line.releases)) return;
      candidates.push(line.current.url);
      line.releases.forEach(function (release) {
        if (release && typeof release.url === "string") candidates.push(release.url);
      });
    });
    return candidates.map(function (url) {
      return { url: url, path: new URL(url, window.location.href).pathname.replace(/\/$/, "") + "/" };
    }).filter(function (candidate) {
      return window.location.pathname.indexOf(candidate.path) === 0;
    }).sort(function (a, b) {
      return b.path.length - a.path.length;
    })[0];
  }

  function isActive(base, active) {
    return active && new URL(base, window.location.href).pathname.replace(/\/$/, "") + "/" === active.path;
  }

  function item(label, base, active) {
    var li = document.createElement("li");
    var link = document.createElement("a");
    link.href = new URL(base, window.location.href).href;
    link.textContent = label;
    if (active) link.setAttribute("aria-current", "page");
    li.appendChild(link);
    return li;
  }

  function appendLine(fragment, line, active) {
    if (!line || typeof line.line !== "string" || !line.current || typeof line.current.url !== "string" || !Array.isArray(line.releases)) {
      throw new Error("invalid release line");
    }
    var headingItem = document.createElement("li");
    var heading = document.createElement("h2");
    heading.textContent = line.line;
    headingItem.appendChild(heading);
    fragment.appendChild(headingItem);
    var currentLabel = line.line === "v1.x" ? "v1.x preview" : "Current: v0.x";
    fragment.appendChild(item(currentLabel, line.current.url, isActive(line.current.url, active)));
    line.releases.forEach(function (release) {
      if (!release || typeof release.url !== "string" || typeof release.version !== "string") {
        throw new Error("invalid release entry");
      }
      fragment.appendChild(item(release.version, release.url, isActive(release.url, active)));
    });
  }

  fetch(menu.dataset.releaseManifest, { credentials: "omit" })
    .then(function (response) {
      if (!response.ok) throw new Error("release manifest unavailable");
      return response.json();
    })
    .then(function (manifest) {
      if (!manifest || !Array.isArray(manifest.lines)) throw new Error("invalid release manifest");
      var fragment = document.createDocumentFragment();
      var active = activeBase(manifest);
      manifest.lines.forEach(function (line) { appendLine(fragment, line, active); });
      menu.replaceChildren(fragment);
    })
    .catch(function () {
      // Keep Hugo-rendered links usable when the manifest is unavailable.
    });
})();
