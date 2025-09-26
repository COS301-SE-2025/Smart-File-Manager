// Modern JavaScript for Smart File Manager Landing Page
document.addEventListener('DOMContentLoaded', function() {
    initializeAnimations();
    initializeFAQ();
    initializeNavigation();
    initializeModals();
    initializeScrollEffects();
});

// Smooth scrolling navigation
function initializeNavigation() {
    const navLinks = document.querySelectorAll('nav a[href^="#"]');

    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();

            const targetId = this.getAttribute('href');
            const targetSection = document.querySelector(targetId);

            if (targetSection) {
                targetSection.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });

    // Mobile menu toggle
    const navToggle = document.querySelector('.nav-toggle');
    const navLinks2 = document.querySelector('.nav-links');

    if (navToggle && navLinks2) {
        navToggle.addEventListener('click', function() {
            navLinks2.classList.toggle('show');
        });
    }
}

// FAQ functionality
function initializeFAQ() {
    const faqItems = document.querySelectorAll('.faq-item');

    faqItems.forEach(item => {
        const question = item.querySelector('.faq-question');

        if (question) {
            question.addEventListener('click', function() {
                // Close other FAQ items
                faqItems.forEach(otherItem => {
                    if (otherItem !== item) {
                        otherItem.classList.remove('active');
                    }
                });

                // Toggle current item
                item.classList.toggle('active');
            });
        }
    });
}

// FAQ toggle function (called from HTML)
function toggleFAQ(element) {
    const faqItem = element.closest('.faq-item');
    const faqItems = document.querySelectorAll('.faq-item');

    // Close other FAQ items
    faqItems.forEach(item => {
        if (item !== faqItem) {
            item.classList.remove('active');
        }
    });

    // Toggle current item
    faqItem.classList.toggle('active');
}

// Modal functionality
function initializeModals() {
    const modals = document.querySelectorAll('.modal');

    // Close modal when clicking outside
    modals.forEach(modal => {
        modal.addEventListener('click', function(e) {
            if (e.target === modal) {
                closeModal();
            }
        });
    });

    // Close modal with Escape key
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            closeModal();
        }
    });
}

// Download functionality
function startDownload(platform) {
    const modal = document.getElementById('download-modal');
    const title = document.getElementById('download-platform-title');

    if (modal && title) {
        title.textContent = `Download for ${platform}`;
        modal.style.display = 'block';
        document.body.style.overflow = 'hidden';

        // Add download tracking
        if (typeof gtag !== 'undefined') {
            gtag('event', 'download_start', {
                'platform': platform
            });
        }

        // Show download instructions based on platform
        updateDownloadInstructions(platform);
    }
}

function updateDownloadInstructions(platform) {
    const instructionsList = document.getElementById('installation-instructions');
    if (!instructionsList) return;

    const instructions = {
        'Windows': [
            'Download the .exe installer file',
            'Right-click and select "Run as administrator"',
            'Follow the installation wizard steps',
            'Launch Smart File Manager from the Start menu'
        ],
        'macOS': [
            'Download the .dmg installer file',
            'Double-click to mount the disk image',
            'Drag Smart File Manager to Applications folder',
            'Launch from Applications or Spotlight search'
        ],
        'Linux': [
            'Download the .deb or .rpm package',
            'Install using your package manager (apt, yum, etc.)',
            'Or run: sudo dpkg -i smart-file-manager.deb',
            'Launch from application menu or terminal'
        ]
    };

    const platformInstructions = instructions[platform] || instructions['Windows'];

    instructionsList.innerHTML = '';
    platformInstructions.forEach(instruction => {
        const li = document.createElement('li');
        li.textContent = instruction;
        instructionsList.appendChild(li);
    });
}

// Close modal
function closeModal() {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        modal.style.display = 'none';
    });
    document.body.style.overflow = 'auto';
}

// Scroll effects
function initializeScrollEffects() {
    // Navbar background on scroll
    const navbar = document.querySelector('nav');

    window.addEventListener('scroll', function() {
        if (window.scrollY > 50) {
            navbar.style.background = 'rgba(46, 46, 46, 0.98)';
            navbar.style.backdropFilter = 'blur(20px)';
        } else {
            navbar.style.background = 'rgba(46, 46, 46, 0.95)';
            navbar.style.backdropFilter = 'blur(20px)';
        }
    });

    // Intersection Observer for animations
    if ('IntersectionObserver' in window) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.style.opacity = '1';
                    entry.target.style.transform = 'translateY(0)';
                }
            });
        }, {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        });

        // Animate elements on scroll
        const animateElements = document.querySelectorAll('.feature-card, .demo-card, .platform-card, .pain-item, .benefit-item');
        animateElements.forEach(el => {
            el.style.opacity = '0';
            el.style.transform = 'translateY(30px)';
            el.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
            observer.observe(el);
        });
    }
}

// Initialize animations
function initializeAnimations() {
    // Hero stats counter animation
    animateCounters();

    // Typing animation for search demo
    animateSearchDemo();
}

// Counter animation for hero stats
function animateCounters() {
    const counters = document.querySelectorAll('.stat-number');

    const animateCounter = (counter) => {
        const target = counter.textContent;
        const isMillions = target.includes('M+');
        const isPercent = target.includes('%+');
        const isSeconds = target.includes('s');

        let finalNumber;
        if (isMillions) {
            finalNumber = parseInt(target.replace('M+', ''));
        } else if (isPercent) {
            finalNumber = parseInt(target.replace('%+', ''));
        } else if (isSeconds) {
            finalNumber = parseInt(target.replace('<', '').replace('s', ''));
        } else {
            finalNumber = parseInt(target);
        }

        if (isNaN(finalNumber)) return;

        let current = 0;
        const increment = finalNumber / 100;
        const timer = setInterval(() => {
            current += increment;

            if (current >= finalNumber) {
                current = finalNumber;
                clearInterval(timer);
            }

            let displayValue = Math.floor(current);
            if (isMillions) {
                counter.textContent = displayValue + 'M+';
            } else if (isPercent) {
                counter.textContent = displayValue + '%+';
            } else if (isSeconds) {
                counter.textContent = '<' + displayValue + 's';
            } else {
                counter.textContent = displayValue;
            }
        }, 20);
    };

    // Start animation when stats come into view
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const counters = entry.target.querySelectorAll('.stat-number');
                counters.forEach(counter => animateCounter(counter));
                observer.unobserve(entry.target);
            }
        });
    });

    const statsSection = document.querySelector('.hero-stats');
    if (statsSection) {
        observer.observe(statsSection);
    }
}

// Typing animation for search demo
function animateSearchDemo() {
    const searchInput = document.querySelector('.search-bar input');
    if (!searchInput) return;

    const phrases = [
        "Find research about machine learning...",
        "budget spreadsheet from last month",
        "presentation about quarterly results",
        "photos from vacation 2023"
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
            // Pause before deleting
            typeSpeed = 2000;
            isDeleting = true;
        } else if (isDeleting && charIndex === 0) {
            isDeleting = false;
            phraseIndex = (phraseIndex + 1) % phrases.length;
            typeSpeed = 500;
        }

        setTimeout(typeWriter, typeSpeed);
    }

    // Start the animation after a delay
    setTimeout(typeWriter, 2000);
}

// Utility functions for backwards compatibility
function showSection(sectionId) {
    const section = document.getElementById(sectionId);
    if (section) {
        section.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
        });
    }
}

function showPrivacyInfo() {
    alert('Smart File Manager is committed to your privacy. All processing happens locally on your device. We never collect, store, or transmit your personal files or data. The application works completely offline and does not require an internet connection.');
}

// Mobile menu toggle
function toggleMobileMenu() {
    const navLinks = document.querySelector('.nav-links');
    if (navLinks) {
        navLinks.classList.toggle('show');
    }
}

// Enhanced download tracking
function trackDownload(platform, method) {
    // Google Analytics tracking
    if (typeof gtag !== 'undefined') {
        gtag('event', 'download_click', {
            'platform': platform,
            'method': method,
            'event_category': 'Downloads',
            'event_label': platform
        });
    }

    // Console log for debugging
    console.log(`Download initiated: ${platform} via ${method}`);
}

// Legacy functions for backwards compatibility
function showLanding() {
    showSection('hero');
}

function showHelp() {
    showSection('faq');
}

function showPopup(platform) {
    startDownload(platform);
}

function showTutorial(tutorial) {
    alert(`The "${tutorial}" tutorial will be available soon!`);
}

function showHelpSection(sectionId) {
    const sections = document.querySelectorAll(".help-section");
    sections.forEach((section) => section.classList.remove("active"));

    const buttons = document.querySelectorAll(".help-nav-btn");
    buttons.forEach((btn) => btn.classList.remove("active"));

    const targetSection = document.getElementById(sectionId);
    if (targetSection) {
        targetSection.classList.add("active");
    }

    if (event && event.target) {
        event.target.classList.add("active");
    }
}

// Performance optimizations
if ('requestIdleCallback' in window) {
    requestIdleCallback(() => {
        // Load non-critical resources
        initializeNonCriticalFeatures();
    });
} else {
    setTimeout(initializeNonCriticalFeatures, 1000);
}

function initializeNonCriticalFeatures() {
    // Initialize features that aren't immediately needed

    // Preload demo videos
    const demoLinks = document.querySelectorAll('.demo-link[href*="drive.google.com"]');
    demoLinks.forEach(link => {
        link.addEventListener('mouseover', function() {
            // Prefetch the video page
            const linkEl = document.createElement('link');
            linkEl.rel = 'prefetch';
            linkEl.href = this.href;
            document.head.appendChild(linkEl);
        });
    });

    // Initialize lazy loading for images if needed
    if ('IntersectionObserver' in window) {
        const images = document.querySelectorAll('img[data-src]');
        const imageObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    img.removeAttribute('data-src');
                    imageObserver.unobserve(img);
                }
            });
        });

        images.forEach(img => imageObserver.observe(img));
    }
}

// Error handling
window.addEventListener('error', function(e) {
    console.error('JavaScript error:', e.error);

    // Optional: Send error to analytics
    if (typeof gtag !== 'undefined') {
        gtag('event', 'exception', {
            'description': e.error.message,
            'fatal': false
        });
    }
});

// Add smooth scrolling for all anchor links
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

// Service Worker registration for offline capabilities (optional)
if ('serviceWorker' in navigator && 'caches' in window) {
    window.addEventListener('load', () => {
        // Uncomment if you want to add service worker
        // navigator.serviceWorker.register('/sw.js')
        //     .then(registration => console.log('SW registered'))
        //     .catch(error => console.log('SW registration failed'));
    });
}