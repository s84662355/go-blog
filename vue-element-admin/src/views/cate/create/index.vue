<template>
	<CateDetail :is-edit="false" />
</template>

<script>
import CateDetail from '../components/CateDetail'

export default {
  name: 'CateDetailCreate',
  components: { CateDetail }
}
</script>

