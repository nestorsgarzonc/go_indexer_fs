<script setup>
import { useEmailsStore } from '@/stores/mails';
import { storeToRefs } from 'pinia';

const emailsStore = useEmailsStore();
const { selectedEmail } = storeToRefs(emailsStore);
</script>

<template>
    <div class="email-preview" v-if="selectedEmail">
        <h1>Selected email: </h1>
        <h2 class="email-preview__subject">{{ selectedEmail.Subject }}</h2>
        <p class="email-preview__info"><b>From:</b> {{ selectedEmail.From }}</p>
        <p class="email-preview__info"><b>To:</b> {{ selectedEmail.To }}</p>
        <p class="email-preview__content">{{ selectedEmail.Content }}</p>
    </div>
</template>



<style scoped>
.email-preview {
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 4px;
}

h1 {
    font-size: 16px;
    font-weight: 600;
}

.email-preview__subject {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 10px;
}

.email-preview__info {
    margin-bottom: 5px;
}

.email-preview__content {
    margin-top: 10px;
}
</style>
