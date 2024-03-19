import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useEmailsStore = defineStore('email', {
  state: () => {
    return {
      emails: [],
      error: null,
      selectedEmail: null,
      isLoading: false,
      query: '',
    }
  },
  getters: {},
  actions: {
    async fetchEmails(query = '') {
      this.isLoading = true
      try {
        const response = await fetch(`http://localhost:8080/emails?query="${query}"`)
        this.emails = await response.json()
      } catch (error) {
        this.error = error
      } finally {
        this.isLoading = false
      }
    },
    async debounceFetchEmails(query) {
      this.isLoading = true
      this.query = query
      await new Promise((resolve) => setTimeout(resolve, 1000))
      if (query === this.query) {
        this.fetchEmails(query)
      }
    },
    setSelectedEmail(email) {
      this.selectedEmail = email
    }
  },
})
