window.addEventListener('load', () => {
    const el = $('#wrapper');
    var handlebar = Handlebars.create();

    // Compile Handlebar Templates
    const errorTemplate = handlebar.compile($('#error-template').html());
    // const headerTemplate = Handlebars.compile($('#header-template').html());
    const todolistTemplate = handlebar.compile($('#todolist-template').html());

    
    // Instantiate api handler
    const api = axios.create({
        baseURL: 'http://localhost:3000',
        timeout: 5000,
    });
    
    // Router Declaration
    const router = new Router({
        mode: 'history',
        page404: (path) => {
            const html = errorTemplate({
                color: 'yellow',
                title: 'Error 404 - Page NOT Found!',
                message: `The path '/${path}' does not exist on this site`,
            });
            el.html(html);
        },
    });
    
    // Display Error Banner
    const showError = (error) => {
        const { title, message } = error;
        const html = errorTemplate({ color: 'red', title, message });
        el.html(html);
    };
    
    router.add('/', async () => {
        let html = todolistTemplate();
        el.html(html);
        try {
            // Load Tasks
            const response = await api.get('/api/tasks');
            const { data } = response;
            
            // Display todolist
            html = todolistTemplate({ data });
            el.html(html);

        } catch (error) {
            showError(error);
        } finally {
                // Remove loader status
            $('.loading').removeClass('loading');
        }
    });
        
    // Navigate app to current url
    router.navigateTo(window.location.pathname);

    // <====== event handler section ========>
});

