class VariantsManager {
  constructor() {
    this.variantsHolder = document.getElementById("variants-holder");
    this.apiUrl = "/api/variants"; // Измените на ваш реальный endpoint
    this.init();
  }

  init() {
    this.loadVariants();
    this.setupEventListeners();
  }

  async loadVariants() {
    try {
      this.showLoading();

      const response = await fetch(this.apiUrl);
      const data = await response.json();

      if (data.success) {
        this.renderVariants(data.variants);
      } else {
        this.showError(data.error || "Ошибка загрузки вариантов");
      }
    } catch (error) {
      console.error("Ошибка загрузки вариантов:", error);
      this.showError("Не удалось загрузить варианты. Проверьте соединение.");
    }
  }

  renderVariants(variants) {
    if (!variants || variants.length === 0) {
      this.showEmptyState();
      return;
    }

    const variantsHTML = variants
      .map((variant) => this.createVariantCard(variant))
      .join("");

    this.variantsHolder.innerHTML = `
            <div class="variants-container">
                ${variantsHTML}
            </div>
        `;

    this.setupCardInteractions();
  }

  createVariantCard(variant) {
    const solvedIcon = variant.solved ? "✅" : "⏳";
    const hasVideo = variant.videoFilePath && variant.videoFilePath !== "";

    return `
            <div class="variant-card" data-uuid="${variant.uuid}">
                <div class="variant-header">
                    <h3 class="variant-title">${this.escapeHtml(variant.name)}</h3>
                    <span class="variant-badge">${variant.class} класс</span>
                </div>

                <div class="variant-meta">
                    <span class="meta-item">
                        <span class="meta-icon">📚</span>
                        ${this.escapeHtml(variant.subject)}
                    </span>
                    <span class="meta-item">
                        <span class="meta-icon">${solvedIcon}</span>
                        ${variant.solved ? "Решен" : "В процессе"}
                    </span>
                </div>

                ${
                  variant.description
                    ? `
                    <p class="variant-description">
                        ${this.escapeHtml(variant.description)}
                    </p>
                `
                    : ""
                }

                <div class="variant-actions">
                    <a href="/api/download/pdf/${variant.uuid}"
                       class="action-btn action-btn-primary"
                       download="${variant.name}.pdf">
                        📄 Скачать PDF
                    </a>

                    ${
                      hasVideo
                        ? `
                        <a href="/api/download/video/${variant.uuid}"
                           class="action-btn action-btn-secondary"
                           download="${variant.name}.mp4">
                            🎥 Смотреть видео
                        </a>
                    `
                        : ""
                    }
                </div>
            </div>
        `;
  }

  setupCardInteractions() {
    const cards = this.variantsHolder.querySelectorAll(".variant-card");

    cards.forEach((card) => {
      card.addEventListener("click", (e) => {
        // Предотвращаем срабатывание при клике на кнопки
        if (!e.target.closest(".action-btn")) {
          const uuid = card.dataset.uuid;
          this.showVariantDetails(uuid);
        }
      });
    });
  }

  showVariantDetails(uuid) {
    // Здесь можно реализовать модальное окно с деталями
    console.log("Показать детали варианта:", uuid);
    // window.open(`/variant/${uuid}`, '_blank');
  }

  showLoading() {
    this.variantsHolder.innerHTML = `
            <div class="loading">
                <div class="loading-spinner"></div>
            </div>
        `;
  }

  showError(message) {
    this.variantsHolder.innerHTML = `
            <div class="error-message">
                <h3>😕 Ошибка загрузки</h3>
                <p>${this.escapeHtml(message)}</p>
                <button onclick="location.reload()" class="action-btn action-btn-primary" style="margin-top: 1rem;">
                    🔄 Попробовать снова
                </button>
            </div>
        `;
  }

  showEmptyState() {
    this.variantsHolder.innerHTML = `
            <div class="empty-state">
                <h3>📝 Варианты не найдены</h3>
                <p>Пока нет доступных вариантов для отображения</p>
            </div>
        `;
  }

  escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
  }

  setupEventListeners() {
    // Обработчик для обновления по кнопке
    document.addEventListener("keydown", (e) => {
      if (e.key === "r" && e.ctrlKey) {
        e.preventDefault();
        this.loadVariants();
      }
    });

    // Авто-обновление каждые 5 минут
    setInterval(
      () => {
        this.loadVariants();
      },
      5 * 60 * 1000,
    );
  }
}

// Инициализация при загрузке страницы
document.addEventListener("DOMContentLoaded", () => {
  new VariantsManager();
});

// Глобальные функции для ручного управления
window.VariantsManager = {
  refresh: function () {
    new VariantsManager().loadVariants();
  },

  filterByClass: function (className) {
    // Можно добавить фильтрацию по классу
    console.log("Фильтр по классу:", className);
  },
};
