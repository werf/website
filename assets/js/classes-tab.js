window.addEventListener('load', () => {
    const wrappers = document.querySelectorAll('[data-tab-group]');

    class GroupTabs {
        constructor(wrapper) {
            this.wrapper = wrapper;
            this.tabs = [];
            this.init = this.initTabs();
            this.currentTab = 1; //null
            this.startActive()
        }

        initTabs() {
            const groupTabs = this;
            const tabs = this.wrapper.querySelectorAll('[data-tabs-item]');

            tabs.forEach((item, idx) => {
                let id = idx + 1;
                const tabContent = item.querySelector('[data-tabs-text]');
                const tabContentHeight = tabContent.offsetHeight;
                const tabDataValue = item.getAttribute('data-tabs-item');
                const tabPicture = this.wrapper.querySelector(`[data-tabs-pic="${tabDataValue}"]`);
                const progressBarContainer = item.querySelector('[data-tabs-progress]');

                this.tabs.push(new Tabs(groupTabs, id, item, tabContent, tabContentHeight, tabPicture, progressBarContainer));
            })
        }

        startActive() {
            const activeTab = this.tabs.find(tab => tab.id === this.currentTab)
            activeTab.setActive()
        }

        removeActive() {
            const activeTab = this.tabs.find(tab => tab.id === this.currentTab)
            activeTab.setInactive();
        }

        switchActiveTab() {
            if (this.currentTab >= this.tabs.length) {
                const activeTab = this.tabs.find(tab => tab.id === 1)
                activeTab.setActive()
            } else {
                const activeTab = this.tabs.find(tab => tab.id === this.currentTab + 1)
                activeTab.setActive()
            }
        }
    }

    class Tabs {
        constructor(groupTab, id, tab, tabContent, tabContentHeight, tabPicture, progressBarContainer) {
            this.id = id;
            this.groupTab = groupTab;
            this.tab = tab;
            this.tabContent = tabContent;
            this.tabContentHeight = tabContentHeight;
            this.tabPicture = tabPicture;
            this.progressBarContainer = progressBarContainer;
            this.progress = null
            this.closeTab();
        }

        closeTab() {
            if (!this.tab.classList.contains('active')) {
                this.tabContent.style.height = '0';
            } else {
                this.tabContent.style.height = this.tabContentHeight + 'px';
            }
        }

        setActive() {
            if (!this.tab.classList.contains('active')) {
                this.groupTab.removeActive();
                this.tab.classList.add('active');
                this.tabContent.style.height = this.tabContentHeight + 'px';
                this.tabPicture.classList.add('active');
                this.groupTab.currentTab = this.id;
                this.setProgressBar(this.progressBarContainer);
            }
        }

        setInactive() {
            this.tab.classList.remove('active');
            this.tabContent.style.height = '0';
            this.tabPicture.classList.remove('active');
            clearInterval(this.progress);
            this.progressBarContainer.innerHTML = '';
        }

        pauseProgress() {
            clearInterval(this.progress);
        }

        playProgress() {
            console.log(123);
        }

        setProgressBar(rootElement) {
            const progressBarNode = document.createElement('span');
            progressBarNode.classList.add('progress-bar');
            rootElement.append(progressBarNode);
            let width = 1;
            const v = this;
            this.progress = setInterval(progressStatus, 10)

            function progressStatus() {
                if (width >= 100) {
                    clearInterval(this.progress)
                    progressBarNode.remove();
                    v.groupTab.switchActiveTab();
                } else {
                    width = width + 0.25;
                    progressBarNode.style.width = width + '%';
                }
            }
        }
    }

    const test = []

    wrappers.forEach((item) => {
        test.push(new GroupTabs(item));
    })

    test[0].tabs.forEach(item => {
        item.tab.addEventListener('click', () => {
            item.setActive();
            // console.log(test);
        })
    })

    test[0].tabs.forEach(item => {
        item.tab.addEventListener('mouseover', () => {
            item.pauseProgress();
        })
    })

    test[0].tabs.forEach(item => {
        item.tab.addEventListener('mouseout', () => {
            item.playProgress();
        })
    })

    test[1].tabs.forEach(item => {
        item.tab.addEventListener('click', () => {
            item.setActive();
        })
    })

    // console.log(test);
});