document.addEventListener("DOMContentLoaded", function () {
  initializeAnimations();
  initializeFAQ();
  initializeNavigation();
  initializeModals();
  initializeScrollEffects();
});

function initializeNavigation() {
  const navLinks = document.querySelectorAll('nav a[href^="#"]');

  navLinks.forEach((link) => {
    link.addEventListener("click", function (e) {
      e.preventDefault();

      const targetId = this.getAttribute("href");
      const targetSection = document.querySelector(targetId);

      if (targetSection) {
        targetSection.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });

        const navLinksContainer = document.querySelector(".nav-links");
        if (navLinksContainer && navLinksContainer.classList.contains("show")) {
          navLinksContainer.classList.remove("show");
        }
      }
    });
  });

  document.addEventListener("click", function (e) {
    const navLinksContainer = document.querySelector(".nav-links");
    const navToggle = document.querySelector(".nav-toggle");
    
    if (
      navLinksContainer &&
      navLinksContainer.classList.contains("show") &&
      !navLinksContainer.contains(e.target) &&
      !navToggle.contains(e.target)
    ) {
      navLinksContainer.classList.remove("show");
    }
  });
}

// FAQ functionality
function initializeFAQ() {
  const faqItems = document.querySelectorAll(".faq-item");

  faqItems.forEach((item) => {
    const question = item.querySelector(".faq-question");

    if (question) {
      question.addEventListener("click", function () {
        faqItems.forEach((otherItem) => {
          if (otherItem !== item) {
            otherItem.classList.remove("active");
          }
        });

        item.classList.toggle("active");
      });
    }
  });
}

function toggleFAQ(element) {
  const faqItem = element.closest(".faq-item");
  const faqItems = document.querySelectorAll(".faq-item");

  // Close other FAQ items
  faqItems.forEach((item) => {
    if (item !== faqItem) {
      item.classList.remove("active");
    }
  });

  // Toggle current item
  faqItem.classList.toggle("active");
}

// Scroll effects
function initializeScrollEffects() {
  // Navbar background on scroll
  const navbar = document.querySelector("nav");

  window.addEventListener("scroll", function () {
    if (window.scrollY > 50) {
      navbar.style.background = "rgba(46, 46, 46, 0.98)";
      navbar.style.backdropFilter = "blur(20px)";
    } else {
      navbar.style.background = "rgba(46, 46, 46, 0.95)";
      navbar.style.backdropFilter = "blur(20px)";
    }
  });

  if ("IntersectionObserver" in window) {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.style.opacity = "1";
            entry.target.style.transform = "translateY(0)";
          }
        });
      },
      {
        threshold: 0.1,
        rootMargin: "0px 0px -50px 0px",
      }
    );

    const animateElements = document.querySelectorAll(
      ".feature-card, .demo-card, .platform-card, .pain-item, .benefit-item"
    );
    animateElements.forEach((el) => {
      el.style.opacity = "0";
      el.style.transform = "translateY(30px)";
      el.style.transition = "opacity 0.6s ease, transform 0.6s ease";
      observer.observe(el);
    });
  }
}

function initializeAnimations() {
  animateSearchDemo();
}

function animateSearchDemo() {
  const searchInput = document.querySelector(".search-bar input");
  if (!searchInput) return;

  const phrases = [
    "Find research about machine learning...",
    "budget spreadsheet from last month",
    "presentation about quarterly results",
    "photos from vacation 2023",
  ];

  let phraseIndex = 0;
  let charIndex = 0;
  let isDeleting = false;

  function typeWriter() {
    const currentPhrase = phrases[phraseIndex];

    if (isDeleting) {
      searchInput.placeholder = currentPhrase.substring(0, charIndex - 1);
      charIndex--;
    } else {
      searchInput.placeholder = currentPhrase.substring(0, charIndex + 1);
      charIndex++;
    }

    let typeSpeed = isDeleting ? 50 : 100;

    if (!isDeleting && charIndex === currentPhrase.length) {
      typeSpeed = 2000;
      isDeleting = true;
    } else if (isDeleting && charIndex === 0) {
      isDeleting = false;
      phraseIndex = (phraseIndex + 1) % phrases.length;
      typeSpeed = 500;
    }

    setTimeout(typeWriter, typeSpeed);
  }

  setTimeout(typeWriter, 2000);
}

function showSection(sectionId) {
  const section = document.getElementById(sectionId);
  if (section) {
    section.scrollIntoView({
      behavior: "smooth",
      block: "start",
    });
  }
}

function showPrivacyInfo() {
  alert(
    "Smart File Manager is committed to your privacy. All processing happens locally on your device. We never collect, store, or transmit your personal files or data. The application works completely offline and does not require an internet connection."
  );
}

function toggleMobileMenu() {
  const navLinks = document.querySelector(".nav-links");
  if (navLinks) {
    navLinks.classList.toggle("show");
  }
}

if ("requestIdleCallback" in window) {
  requestIdleCallback(() => {
    initializeNonCriticalFeatures();
  });
} else {
  setTimeout(initializeNonCriticalFeatures, 1000);
}

function initializeNonCriticalFeatures() {
  // Preload demo videos
  const demoLinks = document.querySelectorAll(
    '.demo-link[href*="drive.google.com"]'
  );
  demoLinks.forEach((link) => {
    link.addEventListener("mouseover", function () {
      // Prefetch the video page
      const linkEl = document.createElement("link");
      linkEl.rel = "prefetch";
      linkEl.href = this.href;
      document.head.appendChild(linkEl);
    });
  });
}

window.addEventListener("error", function (e) {
  console.error("JavaScript error:", e.error);

  if (typeof gtag !== "undefined") {
    gtag("event", "exception", {
      description: e.error.message,
      fatal: false,
    });
  }
});

document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
  anchor.addEventListener("click", function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute("href"));
    if (target) {
      target.scrollIntoView({
        behavior: "smooth",
      });
    }
  });
});

