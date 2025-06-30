// Application State
let bookingState = {
    service: null,
    date: null,
    time: null,
    stylist: null,
    price: 0,
    stylistFee: 0
};

// DOM Elements
const elements = {
    fab: document.getElementById('new-booking-fab'),
    serviceSheet: document.getElementById('service-sheet'),
    datetimeSheet: document.getElementById('datetime-sheet'),
    stylistSheet: document.getElementById('stylist-sheet'),
    confirmSheet: document.getElementById('confirm-sheet'),
    successDialog: document.getElementById('success-dialog')
};

// Initialize the app
document.addEventListener('DOMContentLoaded', function() {
    initializeEventListeners();
    initializeRecentBookings();
});

// Event Listeners
function initializeEventListeners() {
    // FAB click
    elements.fab.addEventListener('click', function() {
        resetBookingState();
        showBottomSheet(elements.serviceSheet);
    });

    // Close buttons
    document.querySelectorAll('.close-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            hideAllBottomSheets();
        });
    });

    // Service selection
    document.querySelectorAll('.service-card').forEach(card => {
        card.addEventListener('click', function() {
            selectService(this);
        });
    });

    // Date selection
    document.querySelectorAll('.date-card').forEach(card => {
        card.addEventListener('click', function() {
            selectDate(this);
        });
    });

    // Time filter chips
    document.querySelectorAll('.filter-chip').forEach(chip => {
        chip.addEventListener('click', function() {
            filterTimes(this);
        });
    });

    // Time selection
    document.querySelectorAll('.time-chip').forEach(chip => {
        chip.addEventListener('click', function() {
            selectTime(this);
        });
    });

    // Stylist shortcuts
    document.querySelectorAll('.shortcut-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            selectStylistShortcut(this);
        });
    });

    // Stylist selection
    document.querySelectorAll('.stylist-card').forEach(card => {
        card.addEventListener('click', function() {
            selectStylist(this);
        });
    });

    // Confirm buttons
    document.getElementById('edit-btn').addEventListener('click', function() {
        hideAllBottomSheets();
        showBottomSheet(elements.serviceSheet);
    });

    document.getElementById('reserve-btn').addEventListener('click', function() {
        makeReservation();
    });

    // Success dialog
    document.getElementById('close-success').addEventListener('click', function() {
        hideSuccessDialog();
        hideAllBottomSheets();
    });

    // Rebook buttons
    document.querySelectorAll('.rebook-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            const card = this.closest('.booking-card');
            rebookAppointment(card);
        });
    });

    // Click outside to close
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('bottom-sheet') && e.target.classList.contains('active')) {
            hideAllBottomSheets();
        }
    });
}

// Service Selection
function selectService(serviceCard) {
    // Remove previous selection
    document.querySelectorAll('.service-card').forEach(card => {
        card.classList.remove('selected');
    });

    // Add selection
    serviceCard.classList.add('selected');

    // Update state
    bookingState.service = {
        id: serviceCard.dataset.service,
        name: serviceCard.querySelector('h3').textContent,
        duration: parseInt(serviceCard.dataset.duration),
        price: parseInt(serviceCard.dataset.price)
    };

    bookingState.price = bookingState.service.price;

    // Auto-advance to next step
    setTimeout(() => {
        hideBottomSheet(elements.serviceSheet);
        showBottomSheet(elements.datetimeSheet);
    }, 300);
}

// Date Selection
function selectDate(dateCard) {
    // Remove previous selection
    document.querySelectorAll('.date-card').forEach(card => {
        card.classList.remove('active');
    });

    // Add selection
    dateCard.classList.add('active');

    // Update state
    bookingState.date = dateCard.dataset.date;

    // Update available times (mock logic)
    updateAvailableTimes();
}

// Time Filtering
function filterTimes(filterChip) {
    // Remove previous selection
    document.querySelectorAll('.filter-chip').forEach(chip => {
        chip.classList.remove('active');
    });

    // Add selection
    filterChip.classList.add('active');

    // Filter time slots
    const filter = filterChip.dataset.filter;
    const timeChips = document.querySelectorAll('.time-chip');

    timeChips.forEach(chip => {
        const time = chip.dataset.time;
        const hour = parseInt(time.split(':')[0]);
        let show = true;

        switch(filter) {
            case 'morning':
                show = hour < 12;
                break;
            case 'afternoon':
                show = hour >= 12 && hour < 17;
                break;
            case 'evening':
                show = hour >= 17;
                break;
            default:
                show = true;
        }

        chip.style.display = show ? 'block' : 'none';
    });
}

// Time Selection
function selectTime(timeChip) {
    // Remove previous selection
    document.querySelectorAll('.time-chip').forEach(chip => {
        chip.classList.remove('selected');
    });

    // Add selection
    timeChip.classList.add('selected');

    // Update state
    bookingState.time = timeChip.dataset.time;

    // Auto-advance to next step
    setTimeout(() => {
        hideBottomSheet(elements.datetimeSheet);
        showBottomSheet(elements.stylistSheet);
    }, 300);
}

// Stylist Shortcuts
function selectStylistShortcut(shortcutBtn) {
    // Remove previous selection
    document.querySelectorAll('.shortcut-btn').forEach(btn => {
        btn.classList.remove('selected');
    });
    document.querySelectorAll('.stylist-card').forEach(card => {
        card.classList.remove('selected');
    });

    // Add selection
    shortcutBtn.classList.add('selected');

    // Update state
    const stylistType = shortcutBtn.dataset.stylist;
    
    if (stylistType === 'none') {
        bookingState.stylist = {
            id: 'none',
            name: '指名なし',
            fee: 0
        };
    } else if (stylistType === 'same') {
        bookingState.stylist = {
            id: 'tanaka',
            name: 'タナカ 美容師（前回と同じ）',
            fee: 0
        };
    }

    bookingState.stylistFee = bookingState.stylist.fee;

    // Auto-advance to confirmation
    setTimeout(() => {
        hideBottomSheet(elements.stylistSheet);
        updateConfirmation();
        showBottomSheet(elements.confirmSheet);
    }, 300);
}

// Stylist Selection
function selectStylist(stylistCard) {
    // Remove previous selection
    document.querySelectorAll('.stylist-card').forEach(card => {
        card.classList.remove('selected');
    });
    document.querySelectorAll('.shortcut-btn').forEach(btn => {
        btn.classList.remove('selected');
    });

    // Add selection
    stylistCard.classList.add('selected');

    // Update state
    bookingState.stylist = {
        id: stylistCard.dataset.stylist,
        name: stylistCard.querySelector('h3').textContent,
        fee: parseInt(stylistCard.dataset.fee)
    };

    bookingState.stylistFee = bookingState.stylist.fee;

    // Auto-advance to confirmation
    setTimeout(() => {
        hideBottomSheet(elements.stylistSheet);
        updateConfirmation();
        showBottomSheet(elements.confirmSheet);
    }, 300);
}

// Update Confirmation
function updateConfirmation() {
    document.getElementById('confirm-service').textContent = bookingState.service.name;
    
    const dateObj = new Date(bookingState.date);
    const formattedDate = `${dateObj.getFullYear()}年${dateObj.getMonth() + 1}月${dateObj.getDate()}日`;
    document.getElementById('confirm-datetime').textContent = `${formattedDate} ${bookingState.time}`;
    
    document.getElementById('confirm-stylist').textContent = bookingState.stylist.name;
    
    const total = bookingState.price + bookingState.stylistFee;
    document.getElementById('confirm-total').textContent = `¥${total.toLocaleString()}`;
}

// Make Reservation
function makeReservation() {
    // Show loading state
    const reserveBtn = document.getElementById('reserve-btn');
    const originalText = reserveBtn.textContent;
    reserveBtn.textContent = '予約中...';
    reserveBtn.disabled = true;

    // Simulate API call
    setTimeout(() => {
        hideBottomSheet(elements.confirmSheet);
        showSuccessDialog();
        
        // Reset button
        reserveBtn.textContent = originalText;
        reserveBtn.disabled = false;

        // Add to recent bookings (mock)
        addToRecentBookings();
    }, 1500);
}

// Rebook Appointment
function rebookAppointment(bookingCard) {
    // Extract booking info (mock data)
    const serviceText = bookingCard.querySelector('h3').textContent;
    const stylistText = bookingCard.querySelector('p').textContent;
    
    // Pre-fill booking state
    resetBookingState();
    
    // Find matching service
    const serviceCards = document.querySelectorAll('.service-card');
    serviceCards.forEach(card => {
        if (card.querySelector('h3').textContent === serviceText) {
            selectService(card);
        }
    });

    // Show date/time selection immediately
    setTimeout(() => {
        hideBottomSheet(elements.serviceSheet);
        showBottomSheet(elements.datetimeSheet);
    }, 100);
}

// Utility Functions
function showBottomSheet(sheet) {
    sheet.classList.add('active');
    document.body.style.overflow = 'hidden';
}

function hideBottomSheet(sheet) {
    sheet.classList.remove('active');
    document.body.style.overflow = '';
}

function hideAllBottomSheets() {
    document.querySelectorAll('.bottom-sheet').forEach(sheet => {
        sheet.classList.remove('active');
    });
    document.body.style.overflow = '';
}

function showSuccessDialog() {
    elements.successDialog.classList.add('active');
}

function hideSuccessDialog() {
    elements.successDialog.classList.remove('active');
}

function resetBookingState() {
    bookingState = {
        service: null,
        date: null,
        time: null,
        stylist: null,
        price: 0,
        stylistFee: 0
    };

    // Reset UI selections
    document.querySelectorAll('.selected, .active').forEach(element => {
        if (!element.classList.contains('date-card') || element.dataset.date !== getTodayString()) {
            element.classList.remove('selected', 'active');
        }
    });

    // Reset date to today
    const todayCard = document.querySelector(`[data-date="${getTodayString()}"]`);
    if (todayCard) {
        todayCard.classList.add('active');
        bookingState.date = getTodayString();
    }

    // Reset time filter
    document.querySelector('[data-filter="all"]').classList.add('active');
    filterTimes(document.querySelector('[data-filter="all"]'));
}

function getTodayString() {
    const today = new Date();
    return today.toISOString().split('T')[0];
}

function updateAvailableTimes() {
    // Mock logic to disable some time slots
    const timeChips = document.querySelectorAll('.time-chip');
    timeChips.forEach((chip, index) => {
        // Randomly disable some slots for demo
        if (Math.random() > 0.7) {
            chip.disabled = true;
            chip.classList.add('disabled');
        } else {
            chip.disabled = false;
            chip.classList.remove('disabled');
        }
    });
}

function addToRecentBookings() {
    // In a real app, this would update the server and refresh the UI
    console.log('Booking added:', bookingState);
}

function initializeRecentBookings() {
    // Initialize with today's date selected
    const todayCard = document.querySelector(`[data-date="${getTodayString()}"]`);
    if (todayCard) {
        todayCard.classList.add('active');
        bookingState.date = getTodayString();
    }

    // Initialize time filters
    document.querySelector('[data-filter="all"]').classList.add('active');
    
    // Update available times
    updateAvailableTimes();
}

// Touch and Gesture Support
let touchStartY = 0;
let touchEndY = 0;

document.addEventListener('touchstart', function(e) {
    touchStartY = e.changedTouches[0].screenY;
});

document.addEventListener('touchend', function(e) {
    touchEndY = e.changedTouches[0].screenY;
    handleSwipe();
});

function handleSwipe() {
    const swipeThreshold = 50;
    const activeSheet = document.querySelector('.bottom-sheet.active');
    
    if (activeSheet && touchStartY - touchEndY < -swipeThreshold) {
        // Swipe down to close
        hideAllBottomSheets();
    }
}

// Keyboard Support
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        hideAllBottomSheets();
        hideSuccessDialog();
    }
});

// Accessibility enhancements
function announceToScreenReader(message) {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.style.position = 'absolute';
    announcement.style.left = '-10000px';
    announcement.textContent = message;
    document.body.appendChild(announcement);
    
    setTimeout(() => {
        document.body.removeChild(announcement);
    }, 1000);
}

// Performance optimization: Lazy load images
function lazyLoadImages() {
    const images = document.querySelectorAll('img[data-src]');
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.classList.remove('lazy');
                imageObserver.unobserve(img);
            }
        });
    });

    images.forEach(img => imageObserver.observe(img));
}

// Initialize lazy loading if images exist
document.addEventListener('DOMContentLoaded', lazyLoadImages);