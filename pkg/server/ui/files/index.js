const selectedFiles = {};
function filesChanged(checked, name) {
    selectionChanged(checked, name, selectedFiles, 'check', disableNavbar, 'files-list-check', getSetAll(getFiles().length));
}

const selectedDestinations = {};
function destinationChanged(checked, name) {
    selectionChanged(checked, name, selectedDestinations, 'check-dest', disableSendDialog, 'check-dest-Destinations', getSetAll(getDestinations().length));
}

function changeAllFiles(checked) {
    changeAll(getFiles(), checked, filesChanged);
}

function changeAllDestinations(checked) {
    changeAll(getDestinations(), checked, destinationChanged);
}

function changeAll(a, checked, change) {
    a.forEach(e => change(checked, e));
}

function selectionChanged(checked, name, set, prefix, disable, all, setAll) {
    delete set[name];
    if (checked) set[name] = true;
    document.getElementById(`${prefix}-${name}`).checked = checked;
    const numChecked = Object.keys(set).length;
    disable(numChecked === 0);
    setAll(document.getElementById(all), numChecked);
}

function getSetAll(total) {
    return (elem, checked) => {
        elem.checked = false;
        elem.indeterminate = false;
        if (checked === total) {
            elem.checked = true;
        } else if (checked > 0) {
            elem.indeterminate = true;
        }
    };
}

function deleteConfirm() {
    const deleteModal = bootstrap.Modal.getOrCreateInstance(document.getElementById('deleteConfirm'));
    setSelectedFiles('delete-files');
    deleteModal.show();
}

function sendConfirm() {
    const sendModal = bootstrap.Modal.getOrCreateInstance(document.getElementById('sendConfirm'));
    setSelectedFiles('send-files');
    uncheckAllDestinations();
    document.getElementById('num-files-badge').innerText = `${Object.keys(selectedFiles).length}`;
    sendModal.show();
}

function setSelectedFiles(id) {
    document.getElementById(id).replaceChildren(...Object.keys(selectedFiles)
        .sort((a, b) => ('' + a).localeCompare(b))
        .map(fn => {
            const elem = document.createElement('li');
            elem.classList.add('list-group-item');
            elem.innerText = fn;
            return elem;
        })
    );
}

function uncheckAll() {
    disableNavbar(true);
    uncheckAllItems(getFiles(), filesChanged, 'files-list-check');
    uncheckAllDestinations();
}

function uncheckAllDestinations() {
    disableSendDialog(true);
    uncheckAllItems(getDestinations(), destinationChanged, 'check-dest-Destinations');
}

function uncheckAllItems(items, change, all) {
    all.checked = false;
    all.indeterminate = false;
    items.forEach(e => change(false, e));
}

function disableNavbar(disable) {
    document.getElementById('delete-button').disabled = disable;
    document.getElementById('send-button').disabled = disable;
}

function disableSendDialog(disable) {
    document.getElementById('send-dialog-button').disabled = disable;
}

function deleteFiles() {
    filesAction('Delete', 'Failed to delete.', 'Deleted files.', 'DELETE', Object.keys(selectedFiles));
}

function sendFiles() {
    filesAction('Send', 'Failed to send.', 'Files sent.', 'POST', {
        files: Object.keys(selectedFiles),
        destinations: Object.keys(selectedDestinations),
    });
}

function filesAction(toastTitle, onfailure, onsuccess, method, body) {
    const toast = (msg, msg2) => prepareToast(toastTitle, `${msg}${!!msg2 ? `: ${msg2}` : ''}`, !!msg2);
    fetch("/", {
        method: method,
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
    }).then(resp => {
        if (resp.status === 200) {
            return toast(onsuccess);
        }
        return resp.text().then(txt => toast(onfailure, txt)).catch(err => toast(onfailure, err));
    }).catch(err => toast(onfailure, err).finally(() => {throw err})).then(() => location.reload());
}

function displayToast(template, container, title, body, isError) {
    const clone = template.content.cloneNode(true);

    const toastTitle = clone.querySelector('#toast-title');
    const toastBody = clone.querySelector('#toast-body');
    toastTitle.innerText = title;
    toastBody.innerText = body;

    if (!isError) {
        const ei = clone.querySelector('#error-icon');
        const parent = clone.querySelector('#header-container');
        parent.removeChild(ei);
    }

    const html = clone.querySelector('#liveToast');
    container.appendChild(html);

    const toast = bootstrap.Toast.getOrCreateInstance(html);
    return new Promise(resolve => resolve())
        .then(() => toast.show())
        .then(() => new Promise(resolve => setTimeout(resolve, 1000*5)))
        .then(() => toast.hide())
        .then(() => new Promise(resolve => setTimeout(resolve, 500)))
        .then(() => container.removeChild(html));
}

function prepareToast(title, body, isError) {
    const toasts = JSON.parse(localStorage.getItem('toasts') ?? '[]');
    toasts.push({title, body, isError});
    localStorage.setItem('toasts', JSON.stringify(toasts));
    return new Promise(resolve => resolve());
}

function displayInitialToasts() {
    const toastTemplate = document.getElementById('toast-template');
    const toastContainer = document.getElementById('toast-container');
    const toasts = JSON.parse(localStorage.getItem('toasts') ?? '[]');
    localStorage.removeItem('toasts');
    toasts.forEach(t => displayToast(toastTemplate, toastContainer, t.title, t.body, t.isError));
}

function openFile(name) {
    window.open(`files/${name}`, "_blank");
}

function init() {
    displayInitialToasts();
    uncheckAll();
}
init();