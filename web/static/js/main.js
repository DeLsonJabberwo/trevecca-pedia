// Dark mode toggle functionality
function toggleDarkMode() {
    if (document.documentElement.classList.contains('dark')) {
        document.documentElement.classList.remove('dark')
        localStorage.theme = 'light'
    } else {
        document.documentElement.classList.add('dark')
        localStorage.theme = 'dark'
    }
}

// Pages drawer toggle functionality
function togglePagesDrawer() {
    const drawer = document.getElementById('pages-drawer')
    const chevron = document.getElementById('drawer-chevron')
    
    if (drawer.classList.contains('hidden')) {
        drawer.classList.remove('hidden')
        chevron.style.transform = 'rotate(180deg)'
    } else {
        drawer.classList.add('hidden')
        chevron.style.transform = 'rotate(0deg)'
    }
}

// Search input focus effects
const searchInput = document.querySelector('input[name="q"]')
if (searchInput) {
    const searchContainer = searchInput.closest('.group')
    if (searchContainer) {
        searchInput.addEventListener('focus', () => {
            searchContainer.classList.add('scale-[1.02]')
        })
        
        searchInput.addEventListener('blur', () => {
            searchContainer.classList.remove('scale-[1.02]')
        })
    }
}

// Smooth scroll for anchor links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault()
        const target = document.querySelector(this.getAttribute('href'))
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            })
        }
    })
})

// Initialize auth state in navigation
updateNavAuth()

// If on profile page, fetch profile data
if (document.getElementById('profile-loading')) {
    fetchProfile()
}

// Wiki editor live preview functionality
document.addEventListener('DOMContentLoaded', function() {
    const editTextarea = document.getElementById('edit-textarea')
    const previewContent = document.getElementById('entry')
    
    if (editTextarea && previewContent) {
        // Debounce function to limit API calls
        let debounceTimer
        function debounce(func, wait) {
            clearTimeout(debounceTimer)
            debounceTimer = setTimeout(func, wait)
        }
        
        // Function to update preview
        async function updatePreview() {
            const content = editTextarea.value
            
            // Don't show preview for empty content
            if (!content || content.trim() === '') {
                previewContent.innerHTML = '<p class="text-gray-400 italic">Start typing to see preview...</p>'
                return
            }
            
            try {
                const response = await fetch('/update-preview', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ content: content })
                })
                
                if (response.ok) {
                    const data = await response.json()
                    previewContent.innerHTML = data.html
                } else {
                    previewContent.innerHTML = '<p class="text-red-500">Failed to generate preview</p>'
                }
            } catch (error) {
                previewContent.innerHTML = '<p class="text-red-500">Error connecting to preview service</p>'
            }
        }
        
        // Initial preview load
        updatePreview()
        
        // Update preview on input with debouncing
        editTextarea.addEventListener('input', function() {
            debounce(updatePreview, 300)
        })
    }
})

