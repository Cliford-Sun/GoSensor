<template>
  <div id="app">
    <div class="container">
      <nav class="sidebar">
        <ul>
          <li><router-link to="/realtime-data">实时数据</router-link></li>
          <li><router-link to="/warning-data">警告数据</router-link></li>
        </ul>
      </nav>
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "App"
};
</script>

<style>
/* 使用Google字体 */
@import url('https://fonts.googleapis.com/css?family=Roboto:400,700&display=swap');

body {
  font-family: 'Roboto', Arial, sans-serif;
  background-color: #f5f5f5;
  margin: 0;
  padding: 0;
}

.container {
  display: flex;
  height: 100vh;
  background-size: cover;
}

.sidebar {
  width: 200px;
  background-color: #e6e6e6;
  padding: 20px;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  width: 100%;
}

.sidebar li {
  margin-bottom: 10px;
}

.sidebar a {
  text-decoration: none;
  color: #333;
  display: block;
  padding: 10px;
  width: 100%;
  transition: all 0.3s ease;
}

.sidebar a:hover {
  transform: scale(1.05);
  background-color: #f0b0b0;
}

.content {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #fff;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  border: 1px solid #ccc;
  border-radius: 5px;
}

@media (max-width: 768px) {
  .sidebar {
    width: 150px;
  }

  .content {
    padding: 10px;
  }
}

.sidebar a:active,
.sidebar a:focus {
  background-color: #d1a1a1;
}
</style>
