import './style.css';
import './app.css';
import * as App from '../wailsjs/go/app/App';

// Application state
class TodoApp {
    constructor() {
        this.tasks = [];
        this.currentEditId = null;
        this.isLoading = false;
        this.currentTab = 'dashboard'; // Отслеживание текущей вкладки
        
        // DOM elements
        this.elements = {};
        
        // Initialize app
        this.init();
    }

    // Initialize application
    async init() {
        this.bindElements();
        this.setupEventListeners();
        this.initTheme();
        
        // Hide loading screen and show app
        await this.loadInitialData();
        this.hideLoading();
    }

    // Bind DOM elements
    bindElements() {
        // Navigation
        this.elements.tabBtns = document.querySelectorAll('.tab-btn');
        this.elements.tabContents = document.querySelectorAll('.tab-content');
        
        // Theme toggle
        this.elements.themeToggle = document.getElementById('themeToggle');
        
        // Dashboard elements
        this.elements.totalTasks = document.getElementById('totalTasks');
        this.elements.pendingTasks = document.getElementById('pendingTasks');
        this.elements.completedTasks = document.getElementById('completedTasks');
        this.elements.overdueTasks = document.getElementById('overdueTasks');
        this.elements.recentTasksList = document.getElementById('recentTasksList');
        this.elements.overdueTasksList = document.getElementById('overdueTasksList');
        
        // Task form
        this.elements.taskForm = document.getElementById('taskForm');
        this.elements.taskTitle = document.getElementById('taskTitle');
        this.elements.taskDescription = document.getElementById('taskDescription');
        this.elements.taskPriority = document.getElementById('taskPriority');
        this.elements.taskDueDate = document.getElementById('taskDueDate');
        this.elements.saveTaskBtn = document.getElementById('saveTaskBtn');
        this.elements.cancelEditBtn = document.getElementById('cancelEditBtn');
        
        // Filters
        this.elements.statusFilter = document.getElementById('statusFilter');
        this.elements.priorityFilter = document.getElementById('priorityFilter');
        this.elements.dateFilter = document.getElementById('dateFilter');
        this.elements.sortBy = document.getElementById('sortBy');
        this.elements.sortOrder = document.getElementById('sortOrder');
        
        // Task lists
        this.elements.activeTasksList = document.getElementById('activeTasksList');
        this.elements.completedTasksList = document.getElementById('completedTasksList');
        this.elements.activeTaskCount = document.getElementById('activeTaskCount');
        this.elements.completedTaskCount = document.getElementById('completedTaskCount');
        
        // Modal
        this.elements.deleteModal = document.getElementById('deleteModal');
        this.elements.deleteTaskTitle = document.getElementById('deleteTaskTitle');
        this.elements.confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
        this.elements.cancelDeleteBtn = document.getElementById('cancelDeleteBtn');
        
        // Quick actions
        this.elements.addTaskBtn = document.getElementById('addTaskBtn');
        this.elements.viewTodayBtn = document.getElementById('viewTodayBtn');
        
        // Toast container
        this.elements.toastContainer = document.getElementById('toastContainer');
    }

    // Setup event listeners
    setupEventListeners() {
        // Tab navigation
        this.elements.tabBtns.forEach(btn => {
            btn.addEventListener('click', (e) => this.switchTab(e.target.dataset.tab));
        });
        
        // Theme toggle
        this.elements.themeToggle.addEventListener('click', () => this.toggleTheme());
        
        // Task form
        this.elements.saveTaskBtn.addEventListener('click', () => this.saveTask());
        this.elements.cancelEditBtn.addEventListener('click', () => this.cancelEdit());
        
        // Filters
        this.elements.statusFilter.addEventListener('change', () => this.filterTasks());
        this.elements.priorityFilter.addEventListener('change', () => this.filterTasks());
        this.elements.dateFilter.addEventListener('change', () => this.filterTasks());
        this.elements.sortBy.addEventListener('change', () => this.filterTasks());
        this.elements.sortOrder.addEventListener('change', () => this.filterTasks());
        
        // Modal
        this.elements.confirmDeleteBtn.addEventListener('click', () => this.confirmDelete());
        this.elements.cancelDeleteBtn.addEventListener('click', () => this.hideDeleteModal());
        
        // Quick actions
        this.elements.addTaskBtn.addEventListener('click', () => this.switchTab('tasks'));
        this.elements.viewTodayBtn.addEventListener('click', () => {
            this.switchTab('tasks');
            this.elements.dateFilter.value = 'today';
            this.filterTasks();
        });
        
        // Collapse buttons
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('collapse-btn')) {
                this.toggleCollapse(e.target);
            }
        });
        
        // Form validation
        this.elements.taskTitle.addEventListener('input', () => this.validateForm());
        
        // Enter key on form
        this.elements.taskTitle.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.saveTask();
            }
        });
    }

    // Initialize theme
    initTheme() {
        const savedTheme = localStorage.getItem('theme') || 'light';
        this.setTheme(savedTheme);
    }

    // Toggle theme
    toggleTheme() {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        this.setTheme(newTheme);
    }

    // Set theme
    setTheme(theme) {
        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('theme', theme);
        
        // Update theme toggle icon
        const icon = this.elements.themeToggle.querySelector('.theme-icon');
        if (theme === 'dark') {
            icon.innerHTML = `
                <circle cx="12" cy="12" r="4"/>
                <path d="M12 1v6M12 17v6M4.22 4.22l4.24 4.24M15.54 15.54l4.24 4.24M1 12h6M17 12h6M4.22 19.78l4.24-4.24M15.54 8.46l4.24-4.24"/>
            `;
        } else {
            icon.innerHTML = `
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            `;
        }
    }

    // Switch tabs with improved data loading
    async switchTab(tabName) {
        this.currentTab = tabName;
        
        // Update tab buttons
        this.elements.tabBtns.forEach(btn => {
            btn.classList.toggle('active', btn.dataset.tab === tabName);
        });
        
        // Update tab contents
        this.elements.tabContents.forEach(content => {
            content.classList.toggle('active', content.id === `${tabName}-tab`);
        });
        
        // Load data for the active tab
        await this.refreshCurrentTabData();
    }

    // Refresh data for current tab
    async refreshCurrentTabData() {
        try {
            if (this.currentTab === 'dashboard') {
                await this.loadDashboardData();
            } else if (this.currentTab === 'tasks') {
                await this.loadTasks();
            }
        } catch (error) {
            console.error('Error refreshing tab data:', error);
        }
    }

    // Hide loading screen
    hideLoading() {
        const loading = document.getElementById('loading');
        const app = document.getElementById('app');
        
        loading.style.display = 'none';
        app.style.display = 'block';
    }

    // Load initial data
    async loadInitialData() {
        try {
            await this.loadDashboardData();
            await this.loadTasks();
        } catch (error) {
            console.error('Error loading initial data:', error);
            this.showToast('Ошибка загрузки данных', 'error');
        }
    }

    // Load dashboard data
    async loadDashboardData() {
        try {
            const data = await App.GetDashboardData();
            this.updateDashboardStats(data.stats);
            this.renderRecentTasks(data.recent_tasks);
            this.renderOverdueTasks(data.overdue_tasks);
        } catch (error) {
            console.error('Error loading dashboard:', error);
            this.showToast('Ошибка загрузки дашборда', 'error');
        }
    }

    // Update dashboard stats
    updateDashboardStats(stats) {
        if (this.elements.totalTasks) {
            this.elements.totalTasks.textContent = stats.total;
            this.elements.pendingTasks.textContent = stats.pending;
            this.elements.completedTasks.textContent = stats.completed;
            this.elements.overdueTasks.textContent = stats.overdue;
        }
    }

    // Render recent tasks
    renderRecentTasks(tasks) {
        const container = this.elements.recentTasksList;
        if (!container) return;
        
        if (!tasks || tasks.length === 0) {
            container.innerHTML = '<div class="empty-state">Нет последних задач</div>';
            return;
        }
        
        container.innerHTML = tasks.slice(0, 5).map(task => `
            <div class="task-item-mini">
                <div class="task-title-mini">${this.escapeHtml(task.title)}</div>
                <div class="task-meta-mini">
                    <span class="priority-${task.priority}">${this.getPriorityLabel(task.priority)}</span>
                    ${task.due_date ? `• ${this.formatDate(task.due_date)}` : ''}
                </div>
            </div>
        `).join('');
    }

    // Render overdue tasks
    renderOverdueTasks(tasks) {
        const container = this.elements.overdueTasksList;
        if (!container) return;
        
        if (!tasks || tasks.length === 0) {
            container.innerHTML = '<div class="empty-state">Нет просроченных задач</div>';
            return;
        }
        
        container.innerHTML = tasks.slice(0, 3).map(task => `
            <div class="task-item-mini">
                <div class="task-title-mini">${this.escapeHtml(task.title)}</div>
                <div class="task-meta-mini">
                    <span class="priority-${task.priority}">${this.getPriorityLabel(task.priority)}</span>
                    • <span style="color: var(--danger);">Просрочено</span>
                </div>
            </div>
        `).join('');
    }

    // Load tasks with better error handling
    async loadTasks() {
        if (this.isLoading) return;
        
        this.setLoading(true);
        
        try {
            const status = this.elements.statusFilter?.value || '';
            const priority = this.elements.priorityFilter?.value || '';
            const sortBy = this.elements.sortBy?.value || 'created_at';
            const sortOrder = this.elements.sortOrder?.value || 'desc';
            
            let tasks;
            
            // Check if date filter is applied
            const dateFilter = this.elements.dateFilter?.value || '';
            if (dateFilter) {
                tasks = await App.GetTasksByDateFilter(dateFilter);
                // Apply additional filters
                if (status) {
                    tasks = tasks.filter(task => task.status === status);
                }
                if (priority) {
                    tasks = tasks.filter(task => task.priority === priority);
                }
            } else {
                tasks = await App.GetTasks(status, priority, sortBy, sortOrder);
            }
            
            this.tasks = tasks || [];
            this.renderTasks();
            
        } catch (error) {
            console.error('Error loading tasks:', error);
            this.showToast('Ошибка загрузки задач', 'error');
        } finally {
            this.setLoading(false);
        }
    }

    // Filter tasks
    async filterTasks() {
        await this.loadTasks();
    }

    // Render tasks with better DOM manipulation
    renderTasks() {
        const activeTasks = this.tasks.filter(task => task.status === 'pending');
        const completedTasks = this.tasks.filter(task => task.status === 'completed');
        
        this.renderTaskList(activeTasks, this.elements.activeTasksList);
        this.renderTaskList(completedTasks, this.elements.completedTasksList);
        
        // Update counters
        if (this.elements.activeTaskCount) {
            this.elements.activeTaskCount.textContent = activeTasks.length;
        }
        if (this.elements.completedTaskCount) {
            this.elements.completedTaskCount.textContent = completedTasks.length;
        }
    }

    // Render task list
    renderTaskList(tasks, container) {
        if (!container) return;
        
        if (tasks.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <svg class="empty-state-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                        <circle cx="12" cy="12" r="10"/>
                        <line x1="12" y1="6" x2="12" y2="10"/>
                        <line x1="12" y1="14" x2="12.01" y2="14"/>
                    </svg>
                    <h3>Нет задач</h3>
                    <p>Задачи в этой категории пока отсутствуют</p>
                </div>
            `;
            return;
        }
        
        container.innerHTML = tasks.map(task => this.renderTaskItem(task)).join('');
    }

    // Render single task item
    renderTaskItem(task) {
        const dueDate = task.due_date ? new Date(task.due_date) : null;
        const isOverdue = task.is_overdue;
        
        return `
            <div class="task-item ${task.status === 'completed' ? 'completed' : ''} ${isOverdue ? 'overdue' : ''}" data-task-id="${task.id}">
                <div class="task-checkbox ${task.status === 'completed' ? 'checked' : ''}" 
                     onclick="todoApp.toggleTaskStatus(${task.id})">
                </div>
                <div class="task-content">
                    <div class="task-title">${this.escapeHtml(task.title)}</div>
                    ${task.description ? `<div class="task-description">${this.escapeHtml(task.description)}</div>` : ''}
                    <div class="task-meta">
                        <span class="task-priority ${task.priority}">${this.getPriorityLabel(task.priority)}</span>
                        ${dueDate ? `<span class="task-due-date ${isOverdue ? 'overdue' : ''}">${this.formatDate(dueDate)}</span>` : ''}
                        <span class="task-created">Создано: ${this.formatDate(new Date(task.created_at))}</span>
                    </div>
                </div>
                <div class="task-actions">
                    <button class="task-action-btn edit" onclick="todoApp.editTask(${task.id})" title="Редактировать">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                    </button>
                    <button class="task-action-btn delete" onclick="todoApp.showDeleteModal(${task.id})" title="Удалить">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                            <polyline points="3,6 5,6 21,6"/>
                            <path d="M19,6v14a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6m3,0V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2v2"/>
                        </svg>
                    </button>
                </div>
            </div>
        `;
    }

    // Save task with comprehensive updates
    async saveTask() {
        const title = this.elements.taskTitle.value.trim();
        const description = this.elements.taskDescription.value.trim();
        const priority = this.elements.taskPriority.value;
        const dueDate = this.elements.taskDueDate.value;
        
        // Validation
        if (!title) {
            this.showToast('Название задачи обязательно', 'error');
            return;
        }
        
        if (title.length > 255) {
            this.showToast('Название задачи слишком длинное', 'error');
            return;
        }
        
        if (description.length > 1000) {
            this.showToast('Описание задачи слишком длинное', 'error');
            return;
        }
        
        this.setLoading(true);
        
        try {
            let result;
            
            if (this.currentEditId) {
                // Update existing task
                result = await App.UpdateTask(
                    this.currentEditId,
                    title,
                    description,
                    '', // status - don't change
                    priority,
                    dueDate ? new Date(dueDate).toISOString() : ''
                );
                this.showToast('Задача обновлена', 'success');
            } else {
                // Create new task
                result = await App.CreateTask(
                    title,
                    description,
                    priority,
                    dueDate ? new Date(dueDate).toISOString() : ''
                );
                this.showToast('Задача создана', 'success');
            }
            
            this.clearForm();
            
            // Comprehensive data reload
            await this.refreshAllData();
            
        } catch (error) {
            console.error('Error saving task:', error);
            this.showToast('Ошибка сохранения задачи', 'error');
        } finally {
            this.setLoading(false);
        }
    }

    // Refresh all data comprehensively
    async refreshAllData() {
        try {
            // Always reload both dashboard and tasks data
            await Promise.all([
                this.loadDashboardData(),
                this.loadTasks()
            ]);
            
            // Small delay to ensure DOM updates
            await new Promise(resolve => setTimeout(resolve, 100));
            
        } catch (error) {
            console.error('Error refreshing all data:', error);
        }
    }

    // Edit task
    editTask(taskId) {
        const task = this.tasks.find(t => t.id === taskId);
        if (!task) return;
        
        this.currentEditId = taskId;
        
        // Fill form with task data
        this.elements.taskTitle.value = task.title;
        this.elements.taskDescription.value = task.description || '';
        this.elements.taskPriority.value = task.priority;
        
        if (task.due_date) {
            const date = new Date(task.due_date);
            this.elements.taskDueDate.value = this.formatDateTimeLocal(date);
        } else {
            this.elements.taskDueDate.value = '';
        }
        
        // Update form UI
        this.elements.saveTaskBtn.innerHTML = `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            Обновить задачу
        `;
        this.elements.cancelEditBtn.style.display = 'inline-flex';
        
        // Scroll to form
        this.elements.taskForm.scrollIntoView({ behavior: 'smooth' });
    }

    // Cancel edit
    cancelEdit() {
        this.currentEditId = null;
        this.clearForm();
    }

    // Clear form
    clearForm() {
        this.elements.taskTitle.value = '';
        this.elements.taskDescription.value = '';
        this.elements.taskPriority.value = 'medium';
        this.elements.taskDueDate.value = '';
        
        this.elements.saveTaskBtn.innerHTML = `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            Добавить задачу
        `;
        this.elements.cancelEditBtn.style.display = 'none';
        this.currentEditId = null;
    }

    // Toggle task status with comprehensive updates
    async toggleTaskStatus(taskId) {
        this.setLoading(true);
        
        try {
            await App.ToggleTaskComplete(taskId);
            this.showToast('Статус задачи изменен', 'success');
            
            // Comprehensive data reload
            await this.refreshAllData();
            
        } catch (error) {
            console.error('Error toggling task status:', error);
            this.showToast('Ошибка изменения статуса', 'error');
        } finally {
            this.setLoading(false);
        }
    }

    // Show delete modal
    showDeleteModal(taskId) {
        const task = this.tasks.find(t => t.id === taskId);
        if (!task) return;
        
        this.currentDeleteId = taskId;
        this.elements.deleteTaskTitle.textContent = task.title;
        this.elements.deleteModal.classList.add('show');
    }

    // Hide delete modal
    hideDeleteModal() {
        this.elements.deleteModal.classList.remove('show');
        this.currentDeleteId = null;
    }

    // Confirm delete with comprehensive updates
    async confirmDelete() {
        if (!this.currentDeleteId) return;
        
        this.setLoading(true);
        
        try {
            await App.DeleteTask(this.currentDeleteId);
            this.showToast('Задача удалена', 'success');
            
            this.hideDeleteModal();
            
            // Comprehensive data reload
            await this.refreshAllData();
            
            // Show warning about UI update issue
            this.showUpdateWarning();
            
        } catch (error) {
            console.error('Error deleting task:', error);
            this.showToast('Ошибка удаления задачи', 'error');
        } finally {
            this.setLoading(false);
        }
    }

    // Show warning about UI update issue (temporary solution)
    showUpdateWarning() {
        // Only show this warning occasionally, not every time
        if (Math.random() < 0.3) { // 30% chance to show
            this.showToast(
                'Если изменения не отображаются, переключитесь между вкладками "Дашборд" и "Задачи" для обновления интерфейса', 
                'warning', 
                'Временная проблема с обновлением', 
                7000 // Show for 7 seconds
            );
        }
    }

    // Show initial warning modal on app startup
    showInitialWarningModal() {
        // Check if warning was already shown in this session
        if (sessionStorage.getItem('warningShown')) {
            return;
        }
        
        // Create warning modal
        const warningModal = document.createElement('div');
        warningModal.className = 'modal show';
        warningModal.id = 'warningModal';
        warningModal.innerHTML = `
            <div class="modal-content" style="max-width: 500px;">
                <div class="modal-header">
                    <h3 style="color: var(--warning); display: flex; align-items: center; gap: 0.5rem;">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" style="width: 24px; height: 24px;">
                            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                            <line x1="12" y1="9" x2="12" y2="13"/>
                            <line x1="12" y1="17" x2="12.01" y2="17"/>
                        </svg>
                        Важное уведомление
                    </h3>
                </div>
                <div class="modal-body">
                    <p style="margin-bottom: 1rem; color: var(--text);">
                        <strong>Внимание!</strong> В приложении есть известная проблема с обновлением интерфейса.
                    </p>
                    <p style="margin-bottom: 1rem; color: var(--text-secondary);">
                        После выполнения операций (создание, редактирование, удаление задач) 
                        изменения могут не отображаться сразу.
                    </p>
                    <div style="background: var(--bg-secondary); padding: 1rem; border-radius: var(--radius); border-left: 4px solid var(--warning); margin-bottom: 1rem;">
                        <p style="color: var(--text); font-weight: 600; margin-bottom: 0.5rem;">
                            Как исправить:
                        </p>
                        <p style="color: var(--text-secondary);">
                            Переключайтесь между вкладками <strong>"Дашборд"</strong> и <strong>"Задачи"</strong> 
                            для обновления отображаемых данных.
                        </p>
                    </div>
                    <p style="color: var(--text-secondary); font-size: 0.875rem;">
                        Мы работаем над исправлением этой проблемы.
                    </p>
                </div>
                <div class="modal-footer">
                    <button class="btn btn-primary" id="warningModalOk">Понятно</button>
                </div>
            </div>
        `;
        
        // Add modal to page
        document.body.appendChild(warningModal);
        
        // Add click handler for OK button
      const okButton = warningModal.querySelector('.modal-footer button');


        okButton.addEventListener('click', () => {
            warningModal.remove();
            sessionStorage.setItem('warningShown', 'true');
        });
        
        // Prevent closing by clicking outside
        warningModal.addEventListener('click', (e) => {
            if (e.target === warningModal) {
                // Do nothing - force user to click OK
            }
        });
    }

    // Toggle collapse
    toggleCollapse(button) {
        const targetId = button.dataset.target;
        const target = document.getElementById(targetId);
        
        if (target) {
            const isCollapsed = target.classList.contains('collapsed');
            target.classList.toggle('collapsed', !isCollapsed);
            button.classList.toggle('collapsed', !isCollapsed);
        }
    }

    // Validate form
    validateForm() {
        const title = this.elements.taskTitle.value.trim();
        const isValid = title.length > 0 && title.length <= 255;
        
        this.elements.saveTaskBtn.disabled = !isValid;
        
        if (title.length > 255) {
            this.elements.taskTitle.style.borderColor = 'var(--danger)';
        } else {
            this.elements.taskTitle.style.borderColor = '';
        }
    }

    // Set loading state
    setLoading(loading) {
        this.isLoading = loading;
        
        if (loading) {
            document.body.style.cursor = 'wait';
            if (this.elements.saveTaskBtn) {
                this.elements.saveTaskBtn.disabled = true;
            }
        } else {
            document.body.style.cursor = '';
            this.validateForm();
        }
    }

    // Show toast notification
    showToast(message, type = 'info', title = '', duration = 4000) {
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        
        const icons = {
            success: `<svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                        <path d="M9 12l2 2 4-4"/>
                        <path d="M21 12c0 4.97-4.03 9-9 9s-9-4.03-9-9 4.03-9 9-9 9 4.03 9 9z"/>
                      </svg>`,
            error: `<svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                      <circle cx="12" cy="12" r="10"/>
                      <line x1="15" y1="9" x2="9" y2="15"/>
                      <line x1="9" y1="9" x2="15" y2="15"/>
                    </svg>`,
            warning: `<svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                        <line x1="12" y1="9" x2="12" y2="13"/>
                        <line x1="12" y1="17" x2="12.01" y2="17"/>
                      </svg>`,
            info: `<svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                     <circle cx="12" cy="12" r="10"/>
                     <line x1="12" y1="16" x2="12" y2="12"/>
                     <line x1="12" y1="8" x2="12.01" y2="8"/>
                   </svg>`
        };
        
        toast.innerHTML = `
            ${icons[type] || icons.info}
            <div class="toast-content">
                ${title ? `<div class="toast-title">${title}</div>` : ''}
                <div class="toast-message">${message}</div>
            </div>
            <button class="toast-close" onclick="this.parentElement.remove()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <line x1="18" y1="6" x2="6" y2="18"/>
                    <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
            </button>
        `;
        
        this.elements.toastContainer.appendChild(toast);
        
        // Auto remove after duration
        setTimeout(() => {
            if (toast.parentElement) {
                toast.remove();
            }
        }, duration);
    }

    // Utility functions
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    getPriorityLabel(priority) {
        const labels = {
            high: 'Высокий',
            medium: 'Средний',
            low: 'Низкий'
        };
        return labels[priority] || priority;
    }

    formatDate(date) {
        if (!date) return '';
        
        const d = new Date(date);
        const now = new Date();
        const diff = d.getTime() - now.getTime();
        const days = Math.ceil(diff / (1000 * 60 * 60 * 24));
        
        if (days === 0) {
            return 'Сегодня';
        } else if (days === 1) {
            return 'Завтра';
        } else if (days === -1) {
            return 'Вчера';
        } else if (days > 0 && days <= 7) {
            return `Через ${days} дн.`;
        } else if (days < 0 && days >= -7) {
            return `${Math.abs(days)} дн. назад`;
        }
        
        return d.toLocaleDateString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });
    }

    formatDateTimeLocal(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        
        return `${year}-${month}-${day}T${hours}:${minutes}`;
    }

    // Debounce utility
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
}

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.todoApp = new TodoApp();
});

// Handle modal clicks outside content
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('modal')) {
        e.target.classList.remove('show');
    }
});

// Handle escape key for modals
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        const openModal = document.querySelector('.modal.show');
        if (openModal) {
            openModal.classList.remove('show');
        }
    }
});