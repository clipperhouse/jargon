using System.Collections.Generic;
using System.IO;
using Xunit;
using System.Linq;
using Jargon.Impl;
using System;

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
            got.AddRange(Jargon.Lemmatize(original, lem));
            
            var expected = new List<Token>();
            expected.AddRange(Jargon.Tokenize("Here is the story of ruby-on-rails node.js, \"javascript\", html5 and asp.net-mvc plus tcpip."));

            Assert.Equal(expected.Count, got.Count);
            for(var i = 0; i < got.Count; i++)
            {
                Assert.Equal(expected[i].String, got[i].String);
            }

            var lemmas = new[] { "ruby-on-rails", "node.js", "javascript", "html5", "asp.net-mvc" };

            var lookup = new Dictionary<string, Token>();
            foreach(var g in got)
            {
                lookup[g.String] = g;
            }

            foreach(var lemma in lemmas)
            {
                var matching = got.Where(g => g.String == lemma).ToList();
                Assert.True(matching.Count > 0);

                var ok = lookup.TryGetValue(lemma, out var l);
                Assert.True(ok);
                Assert.True(l.IsLemma);
            }
        }

        [Fact]
        public void CSV()
        {
            var dict = Data.StackExchange.Instance;
            var lem = new Lemmatizer(dict, 3);

            var original = "\"Ruby on Rails\", 3.4, \"foo\"\n\"bar\",42, \"java script\"";

            var got = new List<Token>();
            got.AddRange(lem.Lemmatize(original));

            var expected = new List<Token>();
            expected.AddRange(Jargon.Tokenize("\"ruby-on-rails\", 3.4, \"foo\"\n\"bar\",42, \"javascript\""));

            Assert.Equal(expected.Count, got.Count);
            for(var i = 0; i < got.Count; i++)
            {
                Assert.Equal(expected[i].String, got[i].String);
            }
        }

        [Fact]
        public void TSV()
        {
            var dict = Data.StackExchange.Instance;
            var lem = new Lemmatizer(dict, 3);
            var original = "Ruby on Rails	3.4	foo\nASPNET	MVC\nbar	42	java script";

            var got = new List<Token>();
            got.AddRange(Jargon.Lemmatize(original, lem));

            var expected = new List<Token>();
            expected.AddRange(Jargon.Tokenize("ruby-on-rails	3.4	foo\nasp.net	model-view-controller\nbar	42	javascript"));

            Assert.Equal(expected.Count, got.Count);
            for (var i = 0; i < got.Count; i++)
            {
                Assert.Equal(expected[i].String, got[i].String);
            }
        }

        class _WordRun : ILemmatizingDictionary
        {
            public (string Canonical, bool Found) Lookup(string[] term, int termLen)
            {
                throw new System.NotImplementedException();
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

            var fakeLem = new Lemmatizer(new _WordRun(), 4);

            using (var r = new StringReader(original))
            using (var tokens = new TextTokens(r))
            using (var sc = new LemmaTokens<TextTokens>(fakeLem, tokens))
            {
                foreach(var kv in expecteds)
                {
                    var take = kv.Key;
                    var expected = kv.Value;

                    var ok = sc.WordRun(take, out var consumed);
                    var taken = sc.WordRunBuffer.AsArray;
                    var takenLen = sc.WordRunBuffer.Count;
                    string[] takenStrsArr;

                    if (ok)
                    {
                        var takenStrs = new List<string>();
                        for (var i = 0; i < takenLen; i++)
                        {
                            var t = taken[i];
                            takenStrs.Add(t);
                        }
                        takenStrsArr = takenStrs.ToArray();
                    }
                    else
                    {
                        consumed = 0;
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
