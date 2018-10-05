using Xunit;

namespace Jargon.Tests
{
    public class DictionaryTests
    {
        [Theory]
        [InlineData("foo.js", "foojs")]
        [InlineData(".Net", ".net")]
        [InlineData("ASP.net-mvc", "aspnetmvc")]
        [InlineData("os/2", "os2")]
        public void TestNormalize_StackExchange(string given, string expected)
        {
            var got = Data.StackExchange.Normalize(given);
            Assert.Equal(expected, got);
        }

        [Theory]
        [InlineData(new[] { "three" }, "3", true)]
        [InlineData(new[] { "five" }, "5", true)]
        [InlineData(new[] { "thirtyfive" }, "35", true)]
        [InlineData(new[] { "thirty-five" }, "35", true)]
        [InlineData(new[] { "three", "hundred" }, "300", true)]
        [InlineData(new[] { "3", "hundred" }, "300", true)]
        [InlineData(new[] { "+3", "hundred" }, "300", true)]
        [InlineData(new[] { "-5", "billion" }, "-5000000000", true)]
        [InlineData(new[] { "3", "hundred", "million" }, "300000000", true)]
        [InlineData(new[] { "0.25" }, "0.25", true)]
        [InlineData(new[] { "4.58", "hundred" }, "458", true)]
        [InlineData(new[] { "4.581", "hundred" }, "458.1", true)]
        [InlineData(new[] { "02134" }, "", false)]
        [InlineData(new[] { "2134" }, "2134", true)]
        [InlineData(new[] { "+011" }, "", false)]
        [InlineData(new[] { "-023" }, "", false)]
        [InlineData(new[] { "foo" }, "", false)]
        [InlineData(new[] { "foo three" }, "", false)]
        [InlineData(new[] { "foo 3" }, "", false)]
        [InlineData(new[] { "hundred" }, "", false)]
        [InlineData(new[] { "hundred", "3" }, "", false)]
        [InlineData(new[] { "million", "seven" }, "", false)]
        [InlineData(new[] { "three", "foo" }, "", false)]
        [InlineData(new[] { "3", "foo" }, "", false)]
        [InlineData(new[] { "three", "hundred", "foo" }, "", false)]
        [InlineData(new[] { "a", "hundred" }, "", false)]
        public void TestInts_Numbers(string[] given, string expectedCanonical, bool expectedFound)
        {
            var got = Data.Numbers.Instance.Lookup(given, given.Length);
            Assert.Equal(expectedCanonical, got.Canonical);
            Assert.Equal(expectedFound, got.Found);
        }
    }
}
