using System.Collections.Generic;
using System.IO;
using Xunit;
using System.Linq;

namespace Jargon.Tests
{
    public class LemmaTokenTests
    {
        [Fact]
        public void Lemmatize()
        {
            var dict = Data.StackExchange.Instance;
            var lem = new Lemmatizer(dict, 3);

            var got = new List<Token>();
            var original = "Here is the story of Ruby on Rails nodeJS, \"Java Script\", html5 and ASPNET mvc plus TCP/IP.";
            using (var r1 = new StringReader(original))
            using (var tokens = new TextTokens(r1))
            using (var l = new LemmaTokens(in lem, tokens))
            {
                while (l.MoveNext())
                {
                    got.Add(l.Current);
                }
            }

            var expected = new List<Token>();
            using (var r2 = new StringReader("Here is the story of ruby-on-rails node.js, \"javascript\", html5 and asp.net-mvc plus tcpip."))
            using (var tokens = new TextTokens(r2))
            {
                while (tokens.MoveNext())
                {
                    expected.Add(tokens.Current);
                }
            }

            Assert.Equal(expected.Count, got.Count);
            for(var i = 0; i < got.Count; i++)
            {
                Assert.Equal(expected[i].Value, got[i].Value);
            }

            var lemmas = new[] { "ruby-on-rails", "node.js", "javascript", "html5", "asp.net-mvc" };

            var lookup = new Dictionary<string, Token>();
            foreach(var g in got)
            {
                lookup[g.Value] = g;
            }

            foreach(var lemma in lemmas)
            {
                var matching = got.Where(g => g.Value == lemma).ToList();
                Assert.True(matching.Count > 0);

                var ok = lookup.TryGetValue(lemma, out var l);
                Assert.True(ok);
                Assert.True(l.Lemma);
            }
        }

        [Fact]
        public void CSV()
        {
            var dict = Data.StackExchange.Instance;
            var lem = new Lemmatizer(dict, 3);

            var original = "\"Ruby on Rails\", 3.4, \"foo\"\n\"bar\",42, \"java script\"";

            var got = new List<Token>();
            using (var r1 = new StringReader(original))
            using (var tokens = new TextTokens(r1))
            using (var l = new LemmaTokens(in lem, tokens))
            {
                while (l.MoveNext())
                {
                    got.Add(l.Current);
                }
            }

            var expected = new List<Token>();
            using (var r2 = new StringReader("\"ruby-on-rails\", 3.4, \"foo\"\n\"bar\",42, \"javascript\""))
            using (var tokens = new TextTokens(r2))
            {
                while (tokens.MoveNext())
                {
                    expected.Add(tokens.Current);
                }
            }

            Assert.Equal(expected.Count, got.Count);
            for(var i = 0; i < got.Count; i++)
            {
                Assert.Equal(expected[i].Value, got[i].Value);
            }
        }

        [Fact]
        public void TSV()
        {
            var dict = Data.StackExchange.Instance;
            var lem = new Lemmatizer(dict, 3);
            var original = "Ruby on Rails	3.4	foo\nASPNET	MVC\nbar	42	java script";

            var got = new List<Token>();
            using (var r1 = new StringReader(original))
            using (var tokens = new TextTokens(r1))
            using (var l = new LemmaTokens(in lem, tokens))
            {
                while (l.MoveNext())
                {
                    got.Add(l.Current);
                }
            }

            var expected = new List<Token>();
            using (var r2 = new StringReader("ruby-on-rails	3.4	foo\nasp.net	model-view-controller\nbar	42	javascript"))
            using (var tokens = new TextTokens(r2))
            {
                while (tokens.MoveNext())
                {
                    expected.Add(tokens.Current);
                }
            }

            Assert.Equal(expected.Count, got.Count);
            for (var i = 0; i < got.Count; i++)
            {
                Assert.Equal(expected[i].Value, got[i].Value);
            }
        }

        [Fact]
        public void WordRun()
        {
            var original = "java script and ";
            var expecteds =
                new Dictionary<int, (string[] Taken, int Count, bool Ok)>
                {
                    [4] = (null, 0, false),                             // attempting to get 4 should fail
                    [3] = (new[] { "java", "script", "and" }, 5, true), // attempting to get 3 should work, consuming 5
                    [2] = (new[] { "java", "script" }, 3, true),        // attempting to get 2 should work, consuming 3 tokens (incl the space)
                    [1] = (new[] { "java" }, 1, true),                  // attempting to get 1 should work, and consume only that token
                };

            var @default = default(Lemmatizer);

            using (var r = new StringReader(original))
            using (var tokens = new TextTokens(r))
            using (var sc = new LemmaTokens(in @default, tokens))
            {
                foreach(var kv in expecteds)
                {
                    var take = kv.Key;
                    var expected = kv.Value;

                    var (taken, consumed, ok) = sc.WordRun(take);
                    string[] takenStrsArr;

                    if (taken != null)
                    {
                        var takenStrs = new List<string>();
                        foreach (var t in taken)
                        {
                            takenStrs.Add(t.Value);
                        }
                        takenStrsArr = takenStrs.ToArray();
                    }
                    else
                    {
                        takenStrsArr = null;
                    }

                    var got = (Taken: takenStrsArr, Count: consumed, Ok: ok);

                    Assert.Equal(expected.Ok, got.Ok);
                    Assert.Equal(expected.Count, got.Count);

                    if(expected.Taken == null)
                    {
                        Assert.Null(got.Taken);
                    }
                    else
                    {
                        Assert.Equal(expected.Taken.Length, got.Taken.Length);
                        for(var i = 0; i < expected.Taken.Length; i++)
                        {
                            Assert.Equal(expected.Taken[i], got.Taken[i]);
                        }
                    }
                }
            }
        }
    }
}
