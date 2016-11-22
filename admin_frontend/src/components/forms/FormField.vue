<template>
    <div class="form-group" v-bind:class="{ 'has-error': fieldErrors(field) }">
        <label class="control-label">{{label}}<strong v-if="required" class="text-danger"> *</strong></label>
        <slot></slot>
        <p class="help-block" v-if="help_text">{{help_text}}</p>
        <p class="help-block" v-for="err in fieldErrors(field)">
            {{err.error}}
        </p>
    </div>
</template>
<script>
  export default{
    props: ['errors', 'label', 'field', 'required', 'help_text'],
    methods: {
      fieldErrors (field) {
        if (this.errors && this.errors.fields && this.errors.fields[field]) {
          return this.errors.fields[field]
        }
        return null
      }
    }
  }
</script>
