<script setup lang="ts">
import { ref } from "vue";
import { useForm } from "vee-validate";
import {
  required,
  isValidEmail,
  isValidPhoneNumber,
} from "../../utils/helpers/validation";
import { registerMutation } from "../../utils/constants/queries/auth";
import ErrorIcon from "../../components/icons/Error.vue";
import CloseIcon from "../../components/icons/Close.vue";
import LoadingIcon from "../../components/icons/Loading.vue";

definePageMeta({
  layout: "auth",
});

const rememberMe = ref(false);
const errorMessage = ref("");
const isLoading = ref(false);
const isError = ref(false);

const { defineField, handleSubmit, errors } = useForm({
  validationSchema: {
    firstName: required,
    lastName: required,
    email: isValidEmail,
    phoneNumber: isValidPhoneNumber,
    password: required,
    confirmPassword: required,
  },
});

const [firstName, firstNameProps] = defineField("firstName");
const [lastName, lastNameProps] = defineField("lastName");
const [email, emailProps] = defineField("email");
const [phoneNumber, phoneNumberProps] = defineField("phoneNumber");
const [password, passwordProps] = defineField("password");
const [confirmPassword, confirmPasswordProps] = defineField("confirmPassword");

const { mutate: signupUser } = useMutation(registerMutation);

const onSubmit = handleSubmit(async () => {
  if (password.value !== confirmPassword.value) {
    errorMessage.value = "Passwords do not match!!";
    isError.value = true;
    return;
  }

  isLoading.value = true;
  isError.value = false;

  try {
    const res: any = await signupUser({
      first_name: firstName.value,
      last_name: lastName.value,
      email: email.value,
      phone_number: phoneNumber.value,
      password: password.value,
      remember_me: rememberMe.value,
    });

    if (res?.data?.signup) {
      const token = useCookie("token", {
        maxAge: rememberMe.value ? 60 * 60 * 24 * 7 : 60 * 60 * 24,
        path: "/",
      });

      token.value = res.data.signup.token;
      console.log("✅ Signup success");
      navigateTo("/events");
    }
  } catch (err: any) {
    console.error("Signup Error:", err);
    isError.value = true;
    errorMessage.value =
      err.graphQLErrors?.[0]?.message || err.message || "Signup failed";
  } finally {
    isLoading.value = false;
  }
});

const toggle = () => {
  isError.value = false;
};
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-slate-50 px-4 py-12"
  >
    <!-- Error Toast (Matches Login Style) -->
    <Transition name="fade">
      <div
        v-if="isError"
        class="fixed top-5 right-5 flex items-center w-full max-w-xs p-4 text-slate-500 bg-white rounded-2xl shadow-xl border border-red-100 z-50"
        role="alert"
      >
        <div
          class="inline-flex items-center justify-center flex-shrink-0 w-10 h-10 text-red-500 bg-red-50 rounded-xl"
        >
          <ErrorIcon class="w-6 h-6" />
        </div>
        <div class="ms-3 text-sm font-bold text-slate-700">
          {{ errorMessage }}
        </div>
        <button
          @click="toggle"
          class="ms-auto bg-transparent text-slate-400 hover:text-slate-900 p-1.5"
        >
          <CloseIcon class="w-5 h-5" />
        </button>
      </div>
    </Transition>

    <div class="w-full max-w-2xl">
      <!-- Card Container -->
      <div
        class="bg-white p-8 md:p-12 rounded-[2.5rem] shadow-sm border border-slate-200"
      >
        <div class="flex flex-col items-center mb-10">
          <NuxtLink to="/">
            <img
              src="../../assets/images/logo.png"
              class="h-14 mb-6"
              alt="Minab Logo"
            />
          </NuxtLink>
          <h1
            class="text-3xl font-black text-slate-900 tracking-tight text-center"
          >
            Create Account
          </h1>
          <p class="text-slate-500 text-sm mt-2 text-center">
            Join Minab Events to discover and host local experiences
          </p>
        </div>

        <form class="space-y-6" @submit.prevent="onSubmit">
          <!-- Row 1: Names -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label
                class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
                >First Name</label
              >
              <input
                v-model="firstName"
                v-bind="firstNameProps"
                type="text"
                class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
                placeholder="John"
              />
              <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
                errors.firstName
              }}</span>
            </div>
            <div>
              <label
                class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
                >Last Name</label
              >
              <input
                v-model="lastName"
                v-bind="lastNameProps"
                type="text"
                class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
                placeholder="Doe"
              />
              <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
                errors.lastName
              }}</span>
            </div>
          </div>

          <!-- Row 2: Contact Info -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label
                class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
                >Phone Number</label
              >
              <input
                v-model="phoneNumber"
                v-bind="phoneNumberProps"
                type="text"
                class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
                placeholder="0911..."
              />
              <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
                errors.phoneNumber
              }}</span>
            </div>
            <div>
              <label
                class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
                >Email Address</label
              >
              <input
                v-model="email"
                v-bind="emailProps"
                type="email"
                class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
                placeholder="name@mail.com"
              />
              <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
                errors.email
              }}</span>
            </div>
          </div>

          <!-- Row 3: Passwords -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label
                class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
                >Password</label
              >
              <input
                v-model="password"
                v-bind="passwordProps"
                type="password"
                class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
                placeholder="••••••••"
              />
              <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
                errors.password
              }}</span>
            </div>
            <div>
              <label
                class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
                >Confirm Password</label
              >
              <input
                v-model="confirmPassword"
                v-bind="confirmPasswordProps"
                type="password"
                class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
                placeholder="••••••••"
              />
              <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
                errors.confirmPassword
              }}</span>
            </div>
          </div>

          <!-- Terms & Remember -->
          <div class="flex items-center">
            <input
              id="remember"
              type="checkbox"
              v-model="rememberMe"
              class="w-4 h-4 text-indigo-600 border-slate-300 rounded focus:ring-indigo-500"
            />
            <label for="remember" class="ms-2 text-sm font-bold text-slate-600"
              >Remember me for 7 days</label
            >
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-2xl font-bold text-lg transition-all shadow-lg shadow-indigo-200 disabled:opacity-50 active:scale-95 flex items-center justify-center gap-2"
          >
            <LoadingIcon v-if="isLoading" class="w-5 h-5 animate-spin" />
            {{ isLoading ? "Creating Account..." : "Create Account" }}
          </button>

          <!-- Footer Link -->
          <div class="text-center pt-4 border-t border-slate-50 mt-6">
            <p class="text-sm font-bold text-slate-500">
              Already have an account?
              <NuxtLink
                to="/auth/login"
                class="text-indigo-600 hover:underline ml-1"
                >Sign in</NuxtLink
              >
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
