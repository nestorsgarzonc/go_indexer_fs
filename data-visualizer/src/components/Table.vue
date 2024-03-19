<script setup>
import { useEmailsStore } from '@/stores/mails';
import { storeToRefs } from 'pinia';

const emailsStore = useEmailsStore();
const { emails } = storeToRefs(emailsStore);
const maxLength = 50;

const setSelectedEmail = (email) => {
    emailsStore.setSelectedEmail(email);
};
</script>

<template>
    <div class="flex flex-col">
        <div class="overflow-x-auto sm:mx-0.5 lg:mx-0.5">
            <div class="py-2 inline-block">
                <div class="overflow-hidden">
                    <table class="min-w-full">
                        <thead class="bg-gray-200 border-b">
                            <tr>
                                <th scope="col" class="text-sm font-medium text-gray-900 px-6 py-4 text-left">
                                    Content
                                </th>
                                <th scope="col" class="text-sm font-medium text-gray-900 px-6 py-4 text-left">
                                    Subject
                                </th>
                                <th scope="col" class="text-sm font-medium text-gray-900 px-6 py-4 text-left">
                                    From
                                </th>
                                <th scope="col" class="text-sm font-medium text-gray-900 px-6 py-4 text-left">
                                    To
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="email in emails" :key="email.id" @click="setSelectedEmail(email)"
                                class="bg-white border-b transition duration-300 ease-in-out hover:bg-gray-100 cursor-pointer">
                                <td class="text-sm text-gray-900 font-light px-6 py-4 whitespace-nowrap">
                                    {{ email.Content.slice(0, maxLength) }}
                                </td>
                                <td class="text-sm text-gray-900 font-light px-6 py-4 whitespace-nowrap">
                                    {{ email.Subject.slice(0, maxLength) }}
                                </td>
                                <td class="text-sm text-gray-900 font-light px-6 py-4 whitespace-nowrap">
                                    {{ email.From.slice(0, 20) }}
                                </td>
                                <td class="text-sm text-gray-900 font-light px-6 py-4 whitespace-nowrap">
                                    {{ email.To.slice(0, 20) }}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Add some spacing and alignment */
.container {
    margin: 0 auto;
    max-width: 1200px;
}

/* Style the table */
table {
    width: 100%;
    border-collapse: collapse;
}

th,
td {
    padding: 1rem;
    text-align: left;
}

/* Add hover effect */
tr:hover {
    background-color: #f7fafc;
    cursor: pointer;
}
</style>
