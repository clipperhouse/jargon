using System.Collections.Generic;
using System.IO;
using Xunit;
using System.Linq;

namespace Jargon.Tests
{
    public class HtmlTokensTests
    {
        [Theory]
        [InlineData("<html>\n<p foo=\"bar\">\nHi! Let's talk Ruby on Rails.\n<!-- Ignore ASPNET MVC in comments -->\n</p>\n</html>\n",
             new[]
             {
                "<p foo=\"bar\">", // tags kept whole
		        "\n",            // whitespace preserved
		        "Hi", "!",
                "Ruby", "on", "Rails", // make sure text node got tokenized
		        "<!-- Ignore ASPNET MVC in comments -->", // make sure comment kept whole
		        "</p>",
             }
         )]
        public void Tokenize(string input, string[] expectedTokens)
        {
            var got = new List<Token>();
            using (var reader = new StringReader(input))
            using (var e = new HTMLTokens(reader))
            {
                while (e.MoveNext())
                {
                    got.Add(e.Current);
                }
            }

            foreach(var e in expectedTokens)
            {
                var matching = got.Where(g => g.String == e).ToList();
                Assert.True(matching.Count > 0);
            }
        }
    }
}
