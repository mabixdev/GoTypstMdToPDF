// DOM elements
const form = document.getElementById('conversion-form');
const markdownInput = document.getElementById('markdown-input');
const filenameInput = document.getElementById('filename');
const convertBtn = document.getElementById('convert-btn');
const clearBtn = document.getElementById('clear-btn');
const btnText = document.querySelector('.btn-text');
const btnSpinner = document.querySelector('.btn-spinner');
const statusMessages = document.getElementById('status-messages');
const healthStatus = document.getElementById('health-status');
const statsInfo = document.getElementById('stats-info');
const refreshStatsBtn = document.getElementById('refresh-stats');

// State management
let isConverting = false;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    checkServiceHealth();
    refreshStats();
    setupEventListeners();
});

// Event listeners
function setupEventListeners() {
    form.addEventListener('submit', handleFormSubmit);
    clearBtn.addEventListener('click', clearForm);
    refreshStatsBtn.addEventListener('click', refreshStats);
    
    // Auto-resize textarea
    markdownInput.addEventListener('input', autoResizeTextarea);
}

// Check service health
async function checkServiceHealth() {
    try {
        updateHealthStatus('checking', 'Checking service...');
        
        const response = await fetch('/health');
        const data = await response.json();
        
        if (response.ok && data.status === 'healthy') {
            updateHealthStatus('healthy', '✅ Service is healthy');
        } else {
            updateHealthStatus('unhealthy', `❌ Service unavailable: ${data.error || 'Unknown error'}`);
        }
    } catch (error) {
        updateHealthStatus('unhealthy', `❌ Connection failed: ${error.message}`);
    }
}

// Update health status indicator
function updateHealthStatus(status, message) {
    healthStatus.textContent = message;
    healthStatus.className = `status-${status}`;
}

// Handle form submission
async function handleFormSubmit(event) {
    event.preventDefault();
    
    if (isConverting) return;
    
    const markdownContent = markdownInput.value.trim();
    const filename = filenameInput.value.trim() || 'document.pdf';
    
    // Validation
    if (!markdownContent) {
        showStatusMessage('error', 'Please enter some Markdown content to convert.');
        return;
    }
    
    // Ensure filename ends with .pdf
    const finalFilename = filename.endsWith('.pdf') ? filename : `${filename}.pdf`;
    
    try {
        await convertMarkdownToPDF(markdownContent, { filename: finalFilename });
    } catch (error) {
        console.error('Conversion error:', error);
        showStatusMessage('error', `Conversion failed: ${error.message}`);
    }
}

// Convert Markdown content to PDF
async function convertMarkdownToPDF(markdownContent, options = {}) {
    setConvertingState(true);
    clearStatusMessages();
    
    showStatusMessage('info', 'Starting PDF generation from Markdown...');
    
    try {
        const response = await fetch('/api/convert-to-pdf', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                markdownContent,
                options
            })
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
        }
        
        // Handle PDF download
        const blob = await response.blob();
        
        if (blob.size === 0) {
            throw new Error('Received empty PDF file');
        }
        
        downloadBlob(blob, options.filename || 'document.pdf');
        showStatusMessage('success', `✅ PDF generated successfully! (${formatFileSize(blob.size)})`);
        
        // Refresh stats after successful conversion
        setTimeout(refreshStats, 1000);
        
    } catch (error) {
        throw error;
    } finally {
        setConvertingState(false);
    }
}

// Download blob as file
function downloadBlob(blob, filename) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = filename;
    
    document.body.appendChild(a);
    a.click();
    
    // Cleanup
    setTimeout(() => {
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    }, 100);
}

// Set converting state
function setConvertingState(converting) {
    isConverting = converting;
    convertBtn.disabled = converting;
    
    if (converting) {
        btnText.style.display = 'none';
        btnSpinner.style.display = 'inline';
    } else {
        btnText.style.display = 'inline';
        btnSpinner.style.display = 'none';
    }
}

// Clear form
function clearForm() {
    markdownInput.value = '';
    filenameInput.value = '';
    clearStatusMessages();
    markdownInput.focus();
}

// Show status message
function showStatusMessage(type, message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `status-message status-${type}`;
    messageDiv.textContent = message;
    
    statusMessages.appendChild(messageDiv);
    
    // Auto-remove success/info messages after 5 seconds
    if (type === 'success' || type === 'info') {
        setTimeout(() => {
            if (messageDiv.parentNode) {
                messageDiv.parentNode.removeChild(messageDiv);
            }
        }, 5000);
    }
}

// Clear status messages
function clearStatusMessages() {
    statusMessages.innerHTML = '';
}

// Refresh stats
async function refreshStats() {
    try {
        const response = await fetch('/api/stats');
        const stats = await response.json();
        
        statsInfo.textContent = `Active processes: ${stats.activeProcesses}`;
        
        if (stats.processes && stats.processes.length > 0) {
            const details = stats.processes.map(p => 
                `PID ${p.pid} (${Math.round(p.duration / 1000)}s)`
            ).join(', ');
            statsInfo.textContent += ` [${details}]`;
        }
    } catch (error) {
        console.warn('Failed to refresh stats:', error);
        statsInfo.textContent = 'Stats unavailable';
    }
}

// Auto-resize textarea
function autoResizeTextarea() {
    const textarea = markdownInput;
    textarea.style.height = 'auto';
    textarea.style.height = `${Math.max(400, textarea.scrollHeight)}px`;
}

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Handle browser errors
window.addEventListener('error', (event) => {
    console.error('JavaScript error:', event.error);
    showStatusMessage('error', 'An unexpected error occurred. Please refresh the page and try again.');
});

// Handle unhandled promise rejections
window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason);
    showStatusMessage('error', 'An unexpected error occurred during processing.');
});

// Periodic health checks (every 30 seconds)
setInterval(checkServiceHealth, 30000);

// Periodic stats refresh (every 10 seconds)
setInterval(refreshStats, 10000); 