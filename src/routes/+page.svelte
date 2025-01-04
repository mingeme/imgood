<script lang='ts'>
  import type { FileWithPreview } from '$lib/types';
  import FileUploader from '$lib/components/FileUploader.svelte';

  const { data } = $props();
  const { user } = $derived(data);
  let selectedFiles = $state<FileWithPreview[]>([]);

  function handleFilesSelected(files: FileWithPreview[]) {
    selectedFiles = [...selectedFiles, ...files];
  }

  function clearFiles() {
    selectedFiles.forEach(file => URL.revokeObjectURL(file.preview));
    selectedFiles = [];
  }

  function handleRemove(index: number) {
    URL.revokeObjectURL(selectedFiles[index].preview);
    selectedFiles = selectedFiles.filter((_, i) => i !== index);
  }
</script>

<div class='container is-max-desktop px-4'>
  <div class='image-upload-title'>
    <h1 class='title is-4'>Image Upload</h1>
    <p>5 MB max per file. 10 files max per request.</p>
  </div>

  <FileUploader
    {user}
    files={selectedFiles}
    onFilesSelected={handleFilesSelected}
    onRemove={handleRemove}
    onClear={clearFiles}
  />
</div>

<style>
  .image-upload-title {
    border-bottom: 1px solid;
    margin-block: 1rem;
  }
</style>
