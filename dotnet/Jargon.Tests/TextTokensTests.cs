using System.Collections.Generic;
using System.IO;
using Xunit;
using System.Linq;
using System.Globalization;

namespace Jargon.Tests
{
    public class TextTokensTest
    {
        [Theory]
        [InlineData("Hi! This is a test of tech terms.\nIt should consider F#, C++, .net, Node.JS and 3.141592 and -123 to be their own tokens. \nSimilarly, #hashtag and @handle should work, as should an first.last+@example.com.\nIt should—wait for it—break on things like em-dashes and \"quotes\" and it ends.\nIt'd be great it it’ll handle apostrophes.\n",
            new[]
            {
                "Hi", "!", "a",
                "F#", "C++", ".net", "Node.JS", "3.141592", "-123",
                "#hashtag", "@handle", "first.last+@example.com",
                "should", "—", "wait", "it", "break", "em-dashes", "quotes", "ends",
                "It'd", "it’ll", "apostrophes"
            }
        )]
        public void Tokenize(string input, string[] expectedTokens)
        {
            var got = new List<Token>();
            got.AddRange(Jargon.Tokenize(input));

            foreach(var tok in expectedTokens)
            {
                var matching = got.Where(g => g.String == tok).ToArray();
                Assert.True(matching.Count() > 0);
            }

            var nextToLast = got[got.Count - 2];
            Assert.Equal(".", nextToLast.String);

            var last = got[got.Count - 1];
            Assert.Equal("\n", last.String);

            foreach(var token in got)
            {
                var ee = StringInfo.GetTextElementEnumerator(token.String);
                var numRunes = 0;
                while (ee.MoveNext())
                {
                    numRunes++;
                }

                if (numRunes == 1) continue;

                Assert.False(token.String.EndsWith(",") || token.String.EndsWith("."));
            }
        }

        [Fact]
        public void TestUrls()
        {
            var tests =
                new Dictionary<string, string>
                {
                    ["http://www.google.com"] = "http://www.google.com",                     // as-is
                    ["http://www.google.com."] = "http://www.google.com",                     // "."" should be considered trailing punct
                    ["http://www.google.com/"] = "http://www.google.com/",                    // trailing slash OK
                    ["http://www.google.com/?"] = "http://www.google.com/",                    // "?" should be considered trailing punct
                    ["http://www.google.com/?foo=bar"] = "http://www.google.com/?foo=bar",            // "?" is querystring
                    ["http://www.google.com/?foo=bar."] = "http://www.google.com/?foo=bar",            // trailing "."
                    ["http://www.google.com/?foo=bar&qaz=qux"] = "http://www.google.com/?foo=bar&qaz=qux",    // "?" with &
                    ["http://www.google.com/?foo=bar&qaz=q%20ux"] = "http://www.google.com/?foo=bar&qaz=q%20ux", // with encoding
                    ["//www.google.com"] = "//www.google.com",                          // scheme-relative
                    ["/usr/local/bin/foo.bar"] = "/usr/local/bin/foo.bar",
                    [@"c:\windows\notepad.exe"] = @"c:\windows\notepad.exe",
                };

            foreach (var kv in tests)
            {
                var first = Jargon.Tokenize(kv.Key).First();
                Assert.Equal(kv.Value, first.String);
            }
        }
    }
}
