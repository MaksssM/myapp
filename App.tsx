import React, {useState, useEffect, useRef} from 'react';
import {
  SafeAreaView,
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  StatusBar,
  ActivityIndicator,
  TextInput,
  Animated,
  Image,
  useColorScheme,
  Switch,
} from 'react-native';
import Centrifuge from 'centrifuge';

// Типы данных
interface Post {
  id: number;
  author: string;
  content: string;
  media?: string; // url
  reactions: {[key: string]: number};
  myReaction?: string;
  date: string;
}

const REACTION_TYPES = ['👍', '🔥', '😮', '❤️‍🩹', '🔁', '⭐'];

// URL backend-сервера
const API_URL = 'http://10.0.2.2:8080'; // Для эмулятора Android, для реального устройства — IP вашего ПК

const THEME_COLORS = {
  light: {
    background: '#f8f8f8',
    card: '#fff',
    text: '#222',
    accent: '#007AFF',
    secondary: '#888',
    input: '#f0f0f0',
    border: '#e0e0e0',
  },
  dark: {
    background: '#181A20',
    card: '#23262F',
    text: '#fff',
    accent: '#4F8CFF',
    secondary: '#aaa',
    input: '#23262F',
    border: '#23262F',
  },
};

function useTheme() {
  const scheme = useColorScheme();
  const [isDark, setIsDark] = useState(scheme === 'dark');
  useEffect(() => {
    setIsDark(scheme === 'dark');
  }, [scheme]);
  return {
    isDark,
    setIsDark,
    colors: isDark ? THEME_COLORS.dark : THEME_COLORS.light,
  };
}

const ReactionBar = ({
  post,
  onReact,
}: {
  post: Post;
  onReact: (type: string) => void;
}) => {
  return (
    <View style={{flexDirection: 'row', marginTop: 8}}>
      {REACTION_TYPES.map(type => (
        <TouchableOpacity
          key={type}
          style={{marginRight: 12, opacity: post.myReaction === type ? 1 : 0.7}}
          onPress={() => onReact(type)}>
          <Text style={{fontSize: 20}}>
            {type}{' '}
            <Text style={{fontWeight: 'bold', fontSize: 16}}>
              {post.reactions[type] || 0}
            </Text>
          </Text>
        </TouchableOpacity>
      ))}
    </View>
  );
};

const AnimatedCard = ({
  children,
  index,
}: {
  children: React.ReactNode;
  index: number;
}) => {
  const fadeAnim = useRef(new Animated.Value(0)).current;
  useEffect(() => {
    Animated.timing(fadeAnim, {
      toValue: 1,
      duration: 400 + index * 60,
      useNativeDriver: true,
    }).start();
  }, []);
  return (
    <Animated.View
      style={{
        opacity: fadeAnim,
        transform: [
          {
            scale: fadeAnim.interpolate({
              inputRange: [0, 1],
              outputRange: [0.96, 1],
            }),
          },
        ],
      }}>
      {children}
    </Animated.View>
  );
};

const FeedScreen = ({
  onProfile,
  onPost,
  onViewPost,
  theme,
  posts,
  onRefresh,
}: any) => {
  const [refreshing, setRefreshing] = useState(false);

  const handleRefresh = async () => {
    setRefreshing(true);
    await onRefresh();
    setRefreshing(false);
  };

  return (
    <FlatList
      data={posts}
      keyExtractor={item => item.id.toString()}
      contentContainerStyle={{padding: 16}}
      ListEmptyComponent={
        <Text
          style={{
            color: theme.colors.secondary,
            textAlign: 'center',
            marginTop: 40,
          }}>
          Постов пока нет. Будьте первым!
        </Text>
      }
      ListHeaderComponent={
        <TouchableOpacity
          activeOpacity={0.8}
          style={[styles.createBtn, {backgroundColor: theme.colors.accent}]}
          onPress={onPost}>
          <Text style={{color: '#fff', fontWeight: 'bold', fontSize: 18}}>
            + Новый пост
          </Text>
        </TouchableOpacity>
      }
      renderItem={({item, index}) => (
        <AnimatedCard index={index}>
          <TouchableOpacity
            activeOpacity={0.7}
            onPress={() => onProfile(item.author)}>
            <Text style={[styles.author, {color: theme.colors.accent}]}>
              {' '}
              {item.author}{' '}
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            activeOpacity={0.8}
            onPress={() => onViewPost(item)}>
            <Text style={[styles.content, {color: theme.colors.text}]}>
              {' '}
              {item.content}{' '}
            </Text>
            {item.media && (
              <Image
                source={{uri: item.media}}
                style={styles.media}
                resizeMode="cover"
              />
            )}
          </TouchableOpacity>
          <ReactionBar post={item} onReact={() => {}} />
          <Text style={[styles.date, {color: theme.colors.secondary}]}>
            {' '}
            {item.date}{' '}
          </Text>
        </AnimatedCard>
      )}
      refreshing={refreshing}
      onRefresh={handleRefresh}
    />
  );
};

const PostScreen = ({post, onBack, theme}: any) => (
  <View
    style={[styles.postContainer, {backgroundColor: theme.colors.background}]}>
    <TouchableOpacity onPress={onBack} style={styles.backButton}>
      <Text style={[styles.backButtonText, {color: theme.colors.accent}]}>
        ← Назад
      </Text>
    </TouchableOpacity>
    <Text style={[styles.profileName, {color: theme.colors.accent}]}>
      {post.author}
    </Text>
    <Text
      style={[
        styles.content,
        {color: theme.colors.text, fontSize: 20, marginBottom: 12},
      ]}>
      {post.content}
    </Text>
    {post.media && (
      <Image
        source={{uri: post.media}}
        style={styles.mediaLarge}
        resizeMode="cover"
      />
    )}
    <ReactionBar post={post} onReact={() => {}} />
    <Text style={[styles.date, {color: theme.colors.secondary}]}>
      {post.date}
    </Text>
  </View>
);

const ProfileScreen = ({author, onBack, theme}: any) => (
  <View
    style={[
      styles.profileContainer,
      {backgroundColor: theme.colors.background},
    ]}>
    <TouchableOpacity onPress={onBack} style={styles.backButton}>
      <Text style={[styles.backButtonText, {color: theme.colors.accent}]}>
        ← Назад
      </Text>
    </TouchableOpacity>
    <Text style={[styles.profileName, {color: theme.colors.accent}]}>
      {author}
    </Text>
    <Text style={[styles.profileInfo, {color: theme.colors.secondary}]}>
      Здесь будет профиль пользователя, его посты, подписки и настройки.
    </Text>
  </View>
);

const CreatePostScreen = ({onBack, onCreate, theme}: any) => {
  const [text, setText] = useState('');
  const [media, setMedia] = useState<string | undefined>(undefined);
  return (
    <View
      style={[
        styles.createPostContainer,
        {backgroundColor: theme.colors.background},
      ]}>
      <TouchableOpacity onPress={onBack} style={styles.backButton}>
        <Text style={[styles.backButtonText, {color: theme.colors.accent}]}>
          ← Назад
        </Text>
      </TouchableOpacity>
      <Text style={[styles.profileName, {color: theme.colors.accent}]}>
        Новый пост
      </Text>
      <TextInput
        style={[
          styles.input,
          {backgroundColor: theme.colors.input, color: theme.colors.text},
        ]}
        placeholder="Что нового?"
        placeholderTextColor={theme.colors.secondary}
        value={text}
        onChangeText={setText}
        multiline
      />
      {/* Здесь может быть кнопка для добавления медиа */}
      <TouchableOpacity
        style={[
          styles.createBtn,
          {backgroundColor: theme.colors.accent, marginTop: 24},
        ]}
        activeOpacity={0.8}
        onPress={() => {
          if (text.trim()) {
            onCreate(text, media);
            setText('');
          }
        }}>
        <Text style={{color: '#fff', fontWeight: 'bold', fontSize: 18}}>
          Опубликовать
        </Text>
      </TouchableOpacity>
    </View>
  );
};

const ThemeSwitcher = ({isDark, setIsDark, theme}: any) => (
  <View style={{flexDirection: 'row', alignItems: 'center', marginRight: 12}}>
    <Text style={{color: theme.colors.text, marginRight: 8}}>
      {isDark ? '🌙' : '☀️'}
    </Text>
    <Switch
      value={isDark}
      onValueChange={setIsDark}
      thumbColor={theme.colors.accent}
    />
  </View>
);

const App = () => {
  const theme = useTheme();
  const [screen, setScreen] = useState<'feed' | 'profile' | 'create' | 'post'>(
    'feed',
  );
  const [selectedAuthor, setSelectedAuthor] = useState<string | null>(null);
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  // Загрузка постов с backend
  const fetchPosts = async () => {
    try {
      setLoading(true);
      const res = await fetch(`${API_URL}/posts`);
      const data = await res.json();
      setPosts(data.reverse()); // новые сверху
    } catch (e) {
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  useEffect(() => {
    const centrifuge = new Centrifuge(
      'ws://localhost:8000/connection/websocket',
    );
    centrifuge.on('connect', (ctx: any) => {
      console.log('Connected to Centrifugo:', ctx);
    });
    centrifuge.on('disconnect', (ctx: any) => {
      console.log('Disconnected from Centrifugo:', ctx);
    });
    centrifuge.subscribe('feed', (message: any) => {
      console.log('New message in feed:', message);
    });
    centrifuge.connect();
  }, []);

  const handleProfile = (author: string) => {
    setSelectedAuthor(author);
    setScreen('profile');
  };
  const handleViewPost = (post: Post) => {
    setSelectedPost(post);
    setScreen('post');
  };
  const handleBack = () => {
    setScreen('feed');
    setSelectedAuthor(null);
    setSelectedPost(null);
  };
  const handleCreate = async (text: string, media?: string) => {
    try {
      await fetch(`${API_URL}/posts`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({author: 'Аноним', content: text, media}),
      });
      await fetchPosts();
    } catch (e) {}
    setScreen('feed');
  };

  return (
    <SafeAreaView
      style={[styles.safeArea, {backgroundColor: theme.colors.background}]}>
      <StatusBar
        barStyle={theme.isDark ? 'light-content' : 'dark-content'}
        backgroundColor={theme.colors.background}
      />
      <View style={[styles.header, {backgroundColor: theme.colors.accent}]}>
        <Text style={styles.headerTitle}>Echo</Text>
        <ThemeSwitcher
          isDark={theme.isDark}
          setIsDark={theme.setIsDark}
          theme={theme}
        />
      </View>
      {screen === 'feed' &&
        (loading ? (
          <ActivityIndicator
            style={{marginTop: 50}}
            size="large"
            color={theme.colors.accent}
          />
        ) : (
          <FeedScreen
            onProfile={handleProfile}
            onPost={() => setScreen('create')}
            onViewPost={handleViewPost}
            theme={theme}
            posts={posts}
            onRefresh={fetchPosts}
          />
        ))}
      {screen === 'profile' && selectedAuthor && (
        <ProfileScreen
          author={selectedAuthor}
          onBack={handleBack}
          theme={theme}
        />
      )}
      {screen === 'post' && selectedPost && (
        <PostScreen post={selectedPost} onBack={handleBack} theme={theme} />
      )}
      {screen === 'create' && (
        <CreatePostScreen
          onBack={handleBack}
          onCreate={handleCreate}
          theme={theme}
        />
      )}
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
  },
  header: {
    height: 60,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    elevation: 4,
    marginBottom: 4,
  },
  headerTitle: {
    color: '#fff',
    fontSize: 24,
    fontWeight: 'bold',
    letterSpacing: 1,
  },
  createBtn: {
    borderRadius: 24,
    paddingVertical: 12,
    alignItems: 'center',
    marginBottom: 16,
  },
  card: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 16,
    marginBottom: 18,
    shadowColor: '#000',
    shadowOpacity: 0.08,
    shadowRadius: 8,
    shadowOffset: {width: 0, height: 2},
    elevation: 2,
    borderWidth: 1,
  },
  author: {
    fontWeight: 'bold',
    fontSize: 18,
    marginBottom: 6,
  },
  content: {
    fontSize: 16,
    marginBottom: 8,
  },
  media: {
    width: '100%',
    height: 180,
    borderRadius: 12,
    marginTop: 8,
    marginBottom: 8,
  },
  mediaLarge: {
    width: '100%',
    height: 320,
    borderRadius: 16,
    marginTop: 8,
    marginBottom: 8,
  },
  date: {
    fontSize: 13,
    marginTop: 8,
    alignSelf: 'flex-end',
  },
  profileContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
  },
  profileName: {
    fontSize: 28,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  profileInfo: {
    fontSize: 16,
    marginTop: 24,
    textAlign: 'center',
  },
  backButton: {
    position: 'absolute',
    top: 24,
    left: 16,
    padding: 8,
  },
  backButtonText: {
    fontSize: 18,
  },
  postContainer: {
    flex: 1,
    padding: 24,
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  createPostContainer: {
    flex: 1,
    padding: 24,
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  input: {
    width: '100%',
    minHeight: 80,
    borderRadius: 12,
    padding: 12,
    fontSize: 16,
    marginTop: 12,
    borderWidth: 1,
  },
});

export default App;
