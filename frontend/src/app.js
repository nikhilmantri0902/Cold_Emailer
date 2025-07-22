// Configuration
const API_BASE_URL = '/api';

// Global state
let currentPage = 1;
let currentPageSize = 50;
let currentStage = '';
let allCompanies = [];
let filteredCompanies = [];
let selectedCompanyId = null;
let systemConfig = null;

// Utility functions
function showError(message) {
    const container = document.getElementById('emailLogsContainer');
    container.innerHTML = `<div class="error">❌ Error: ${message}</div>`;
}

function showLoading(containerId) {
    const container = document.getElementById(containerId);
    container.innerHTML = '<div class="loading">⏳ Loading...</div>';
}

// API functions
async function fetchAPI(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// System Config
async function loadSystemConfig() {
    try {
        systemConfig = await fetchAPI('/config');
        renderSystemConfig();
    } catch (error) {
        console.error('Failed to load system config:', error);
    }
}

function renderSystemConfig() {
    if (!systemConfig) return;
    
    const container = document.getElementById('systemConfig');
    container.innerHTML = `
        <div class="filter-group">
            <label>🌍 Target Countries:</label>
            <div style="font-size: 12px; color: #6c757d;">
                ${systemConfig.target_countries.join(', ')}
            </div>
        </div>
        <div class="filter-group">
            <label>👔 Target Roles:</label>
            <div style="font-size: 12px; color: #6c757d;">
                ${systemConfig.suitable_roles.join(', ')}
            </div>
        </div>
        <div class="filter-group">
            <label>🏭 Industries:</label>
            <div style="font-size: 12px; color: #6c757d;">
                ${systemConfig.industries.join(', ')}
            </div>
        </div>
    `;
}

// Email Logs functions
async function loadEmailLogs(page = 1) {
    currentPage = page;
    const stage = document.getElementById('stageFilter').value;
    const pageSize = document.getElementById('pageSizeSelect').value;
    
    currentStage = stage;
    currentPageSize = parseInt(pageSize);
    
    showLoading('emailLogsContainer');
    
    try {
        const params = new URLSearchParams({
            page: page.toString(),
            page_size: pageSize
        });
        
        if (stage) {
            params.append('stage', stage);
        }
        
        const data = await fetchAPI(`/email-logs?${params}`);
        renderEmailLogs(data.logs);
        renderPagination(data.pagination);
    } catch (error) {
        showError(`Failed to load email logs: ${error.message}`);
    }
}

function renderEmailLogs(logs) {
    const container = document.getElementById('emailLogsContainer');
    
    if (!logs || logs.length === 0) {
        container.innerHTML = '<div class="loading">📭 No email logs found</div>';
        return;
    }
    
    const logsHTML = logs.map(log => `
        <div class="log-item">
            <div class="log-header">
                <div class="log-subject">${escapeHtml(log.email_subject || 'No Subject')}</div>
                <div class="log-stage ${log.email_stage}">${log.email_stage}</div>
            </div>
            <div class="log-details">
                <div class="log-detail">
                    <strong>📧 Contact:</strong> ${escapeHtml(log.contact_name || 'Unknown')} 
                    (${escapeHtml(log.contact_email || '')})
                </div>
                <div class="log-detail">
                    <strong>🏢 Company:</strong> ${escapeHtml(log.company_name || 'Unknown')}
                </div>
                <div class="log-detail">
                    <strong>👔 Role:</strong> ${escapeHtml(log.contact_role || 'Unknown')}
                </div>
                <div class="log-detail">
                    <strong>📅 Created:</strong> ${formatDate(log.created_at)}
                </div>
            </div>
            ${log.error_message ? `
                <div class="log-detail" style="margin-top: 10px;">
                    <strong>❌ Error:</strong> <span style="color: #721c24;">${escapeHtml(log.error_message)}</span>
                </div>
            ` : ''}
            ${log.email_body ? `
                <div class="log-body">
                    ${escapeHtml(log.email_body).substring(0, 200)}${log.email_body.length > 200 ? '...' : ''}
                </div>
            ` : ''}
        </div>
    `).join('');
    
    container.innerHTML = `<div class="logs-container">${logsHTML}</div>`;
}

function renderPagination(pagination) {
    const container = document.getElementById('pagination');
    
    if (!pagination || pagination.total_pages <= 1) {
        container.style.display = 'none';
        return;
    }
    
    container.style.display = 'flex';
    
    const buttons = [];
    
    // Previous button
    if (pagination.page > 1) {
        buttons.push(`<button class="page-btn" onclick="loadEmailLogs(${pagination.page - 1})">« Previous</button>`);
    }
    
    // Page numbers
    const startPage = Math.max(1, pagination.page - 2);
    const endPage = Math.min(pagination.total_pages, pagination.page + 2);
    
    for (let i = startPage; i <= endPage; i++) {
        const activeClass = i === pagination.page ? 'active' : '';
        buttons.push(`<button class="page-btn ${activeClass}" onclick="loadEmailLogs(${i})">${i}</button>`);
    }
    
    // Next button
    if (pagination.page < pagination.total_pages) {
        buttons.push(`<button class="page-btn" onclick="loadEmailLogs(${pagination.page + 1})">Next »</button>`);
    }
    
    // Info
    buttons.push(`<span style="margin-left: 20px; color: #6c757d;">
        Page ${pagination.page} of ${pagination.total_pages} (${pagination.total} total)
    </span>`);
    
    container.innerHTML = buttons.join('');
}

// Companies functions
async function loadCompanies() {
    showLoading('companiesContainer');
    
    try {
        const data = await fetchAPI('/companies');
        allCompanies = data.companies || [];
        filteredCompanies = [...allCompanies];
        renderCompanies(filteredCompanies);
    } catch (error) {
        const container = document.getElementById('companiesContainer');
        container.innerHTML = `<div class="error">❌ Failed to load companies: ${error.message}</div>`;
    }
}

function filterCompanies() {
    const searchTerm = document.getElementById('companySearch').value.toLowerCase();
    
    if (!searchTerm.trim()) {
        filteredCompanies = [...allCompanies];
    } else {
        filteredCompanies = allCompanies.filter(company => 
            (company.Name && company.Name.toLowerCase().includes(searchTerm)) ||
            (company.Industry && company.Industry.toLowerCase().includes(searchTerm)) ||
            (company.Website && company.Website.toLowerCase().includes(searchTerm))
        );
    }
    
    renderCompanies(filteredCompanies);
}

function renderCompanies(companies) {
    const container = document.getElementById('companiesContainer');
    
    if (!companies || companies.length === 0) {
        container.innerHTML = '<div class="loading">🏢 No companies found</div>';
        return;
    }
    
    const companiesHTML = companies.map(company => `
        <div class="company-card ${selectedCompanyId === company.ID ? 'selected' : ''}" 
             data-company-id="${company.ID}"
             onclick="selectCompany('${company.ID}', '${escapeHtml(company.Name || 'Unknown Company')}')">
            <div class="company-name">${escapeHtml(company.Name || 'Unknown Company')}</div>
            <div class="company-info">🌐 ${escapeHtml(company.Website || 'No website')}</div>
            <div class="company-info">🏭 ${escapeHtml(company.Industry || 'No industry')}</div>
            <div class="company-info">💻 ${escapeHtml(company.TechDetails || 'No tech details')}</div>
            <div class="company-info">📝 ${escapeHtml((company.CompanyDetails || 'No description').substring(0, 100))}${(company.CompanyDetails || '').length > 100 ? '...' : ''}</div>
        </div>
    `).join('');
    
    container.innerHTML = `<div class="company-list">${companiesHTML}</div>`;
}

async function selectCompany(companyId, companyName) {
    selectedCompanyId = companyId;
    
    // Update visual selection
    document.querySelectorAll('.company-card').forEach(card => {
        card.classList.remove('selected');
    });
    
    // Find and select the clicked company card
    const clickedCard = document.querySelector(`[data-company-id="${companyId}"]`);
    if (clickedCard) {
        clickedCard.classList.add('selected');
    }
    
    // Show contacts section
    const contactsSection = document.getElementById('contactsSection');
    const selectedCompanyNameSpan = document.getElementById('selectedCompanyName');
    const contactsContainer = document.getElementById('contactsContainer');
    
    if (!contactsSection || !selectedCompanyNameSpan || !contactsContainer) {
        console.error('Missing DOM elements for contacts section');
        return;
    }
    
    selectedCompanyNameSpan.textContent = companyName;
    contactsSection.style.display = 'block';
    contactsContainer.innerHTML = '<div class="loading">⏳ Loading contacts...</div>';
    
    try {
        const data = await fetchAPI(`/companies/${companyId}/contacts`);
        console.log('Contacts data received:', data);
        renderContacts(data.contacts);
    } catch (error) {
        console.error('Error loading contacts:', error);
        contactsContainer.innerHTML = `<div class="error">❌ Failed to load contacts: ${error.message}</div>`;
    }
}

function renderContacts(contacts) {
    console.log('renderContacts called with:', contacts);
    const container = document.getElementById('contactsContainer');
    
    if (!container) {
        console.error('contactsContainer element not found');
        return;
    }
    
    if (!contacts || contacts.length === 0) {
        console.log('No contacts found or empty contacts array');
        container.innerHTML = '<div class="loading">👥 No contacts found for this company</div>';
        return;
    }
    
    console.log(`Rendering ${contacts.length} contacts`);
    const contactsHTML = contacts.map(contact => `
        <div class="contact-item">
            <div class="contact-name">${escapeHtml(contact.Name || 'Unknown Contact')}</div>
            <div class="contact-details">
                📧 ${escapeHtml(contact.EmailID || 'No email')} | 
                👔 ${escapeHtml(contact.Role || 'No role')} | 
                📞 ${escapeHtml(contact.PhoneNumber || 'No phone')}
                ${contact.LinkedInURL ? ` | 🔗 <a href="${escapeHtml(contact.LinkedInURL)}" target="_blank">LinkedIn</a>` : ''}
            </div>
        </div>
    `).join('');
    
    container.innerHTML = contactsHTML;
    console.log('Contacts rendered successfully');
}

// Tab switching
function switchTab(tabId) {
    // Update tab buttons
    document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
    event.target.classList.add('active');
    
    // Update tab panels
    document.querySelectorAll('.tab-panel').forEach(panel => panel.classList.remove('active'));
    document.getElementById(tabId).classList.add('active');
    
    // Load data for the selected tab
    if (tabId === 'email-logs') {
        loadEmailLogs();
    } else if (tabId === 'companies') {
        loadCompanies();
    }
}

// Utility functions
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    try {
        return new Date(dateString).toLocaleString();
    } catch {
        return dateString;
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Load system configuration
    loadSystemConfig();
    
    // Load initial data for the active tab
    loadEmailLogs();
    
    // Set up event listeners
    document.getElementById('stageFilter').addEventListener('change', () => loadEmailLogs(1));
    document.getElementById('pageSizeSelect').addEventListener('change', () => loadEmailLogs(1));
    document.getElementById('companySearch').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            filterCompanies();
        }
    });
});

// Expose functions to global scope for onclick handlers
window.switchTab = switchTab;
window.loadEmailLogs = loadEmailLogs;
window.loadCompanies = loadCompanies;
window.filterCompanies = filterCompanies;
window.selectCompany = selectCompany; 