import './App.css'
import Feed from './components/Feed'

function App() {
	return (
		<div className='app-container'>
			<header className='app-header'>
				<h1>TipaTwitter</h1>
			</header>
			<main className='app-main'>
				<Feed />
			</main>
			<footer className='app-footer'>
				<p>© 2025 TipaTwitter. Все права защищены.</p>
			</footer>
		</div>
	)
}

export default App
