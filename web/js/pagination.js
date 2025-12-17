let offset = 0;
const limit = 4;

function loadMore() {
    offset += limit;

    const params = new URLSearchParams(window.location.search);
    params.set("offset", offset);

    fetch("/?" + params.toString())
        .then(response => response.text())
        .then(data => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(data, "text/html");

            const newContentElement = doc.querySelector(".university-list");
            const currentList = document.getElementById("university-list");
            const loadMoreBtn = document.getElementById("load-more");
            const hasMoreElement = doc.getElementById("has-more");

            let hasMore = false;
            if (hasMoreElement) {
                hasMore = hasMoreElement.getAttribute("data-value") === "true";
            }

            if (newContentElement && currentList) {
                currentList.insertAdjacentHTML("beforeend", newContentElement.innerHTML);
            }

            if (!hasMore) {
                loadMoreBtn.style.display = "none";
            }
        })
        .catch(err => {
            console.log("Ошибка при загрузке данных:", err);
            document.getElementById("load-more").style.display = "none";
        });
}
